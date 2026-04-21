package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"easyssl/server/internal/db"
	"easyssl/server/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *db.DB
}

func New(database *db.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) CreateAdmin(ctx context.Context, email, passwordHash string) (*model.Admin, error) {
	id := uuid.NewString()
	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO admins(id,email,password_hash,created_at,updated_at) VALUES($1,$2,$3,$4,$5)`, id, email, passwordHash, now, now)
	if err != nil {
		return nil, err
	}
	return &model.Admin{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: now, UpdatedAt: now}, nil
}

func (r *Repository) GetAdminByEmail(ctx context.Context, email string) (*model.Admin, error) {
	m := &model.Admin{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id,email,password_hash,created_at,updated_at FROM admins WHERE email=$1`, email).
		Scan(&m.ID, &m.Email, &m.PasswordHash, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *Repository) GetAdminByID(ctx context.Context, id string) (*model.Admin, error) {
	m := &model.Admin{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id,email,password_hash,created_at,updated_at FROM admins WHERE id=$1`, id).
		Scan(&m.ID, &m.Email, &m.PasswordHash, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *Repository) UpdateAdminPassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE admins SET password_hash=$2,updated_at=$3 WHERE id=$1`, id, passwordHash, time.Now())
	return err
}

func (r *Repository) ListAccesses(ctx context.Context) ([]model.Access, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id,name,provider,config,reserve,deleted_at,created_at,updated_at FROM accesses WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.Access, 0)
	for rows.Next() {
		var m model.Access
		var raw []byte
		if err := rows.Scan(&m.ID, &m.Name, &m.Provider, &raw, &m.Reserve, &m.DeletedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &m.Config)
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) SaveAccess(ctx context.Context, in model.Access) (*model.Access, error) {
	now := time.Now()
	raw, _ := json.Marshal(in.Config)
	if in.ID == "" {
		in.ID = uuid.NewString()
		in.CreatedAt = now
		in.UpdatedAt = now
		_, err := r.db.Pool.Exec(ctx, `INSERT INTO accesses(id,name,provider,config,reserve,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7)`, in.ID, in.Name, in.Provider, raw, in.Reserve, now, now)
		if err != nil {
			return nil, err
		}
		return &in, nil
	}
	in.UpdatedAt = now
	_, err := r.db.Pool.Exec(ctx, `UPDATE accesses SET name=$2,provider=$3,config=$4,reserve=$5,updated_at=$6 WHERE id=$1`, in.ID, in.Name, in.Provider, raw, in.Reserve, now)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *Repository) SoftDeleteAccess(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE accesses SET deleted_at=$2,updated_at=$2 WHERE id=$1`, id, time.Now())
	return err
}

func (r *Repository) ListWorkflows(ctx context.Context) ([]model.Workflow, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id,name,description,trigger,trigger_cron,enabled,graph_draft,graph_content,has_draft,has_content,last_run_id,last_run_status,last_run_time,created_at,updated_at FROM workflows ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Workflow, 0)
	for rows.Next() {
		var m model.Workflow
		var draftRaw, contentRaw []byte
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Trigger, &m.TriggerCron, &m.Enabled, &draftRaw, &contentRaw, &m.HasDraft, &m.HasContent, &m.LastRunID, &m.LastRunStatus, &m.LastRunTime, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if len(draftRaw) > 0 {
			_ = json.Unmarshal(draftRaw, &m.GraphDraft)
		}
		if len(contentRaw) > 0 {
			_ = json.Unmarshal(contentRaw, &m.GraphContent)
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) GetWorkflow(ctx context.Context, id string) (*model.Workflow, error) {
	m := &model.Workflow{}
	var draftRaw, contentRaw []byte
	err := r.db.Pool.QueryRow(ctx, `SELECT id,name,description,trigger,trigger_cron,enabled,graph_draft,graph_content,has_draft,has_content,last_run_id,last_run_status,last_run_time,created_at,updated_at FROM workflows WHERE id=$1`, id).
		Scan(&m.ID, &m.Name, &m.Description, &m.Trigger, &m.TriggerCron, &m.Enabled, &draftRaw, &contentRaw, &m.HasDraft, &m.HasContent, &m.LastRunID, &m.LastRunStatus, &m.LastRunTime, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(draftRaw) > 0 {
		_ = json.Unmarshal(draftRaw, &m.GraphDraft)
	}
	if len(contentRaw) > 0 {
		_ = json.Unmarshal(contentRaw, &m.GraphContent)
	}
	return m, nil
}

func (r *Repository) SaveWorkflow(ctx context.Context, in model.Workflow) (*model.Workflow, error) {
	now := time.Now()
	draftRaw, _ := json.Marshal(in.GraphDraft)
	contentRaw, _ := json.Marshal(in.GraphContent)
	if in.ID == "" {
		in.ID = uuid.NewString()
		in.CreatedAt = now
		in.UpdatedAt = now
		_, err := r.db.Pool.Exec(ctx, `INSERT INTO workflows(id,name,description,trigger,trigger_cron,enabled,graph_draft,graph_content,has_draft,has_content,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, in.ID, in.Name, in.Description, in.Trigger, in.TriggerCron, in.Enabled, draftRaw, contentRaw, in.HasDraft, in.HasContent, now, now)
		if err != nil {
			return nil, err
		}
		return &in, nil
	}
	in.UpdatedAt = now
	_, err := r.db.Pool.Exec(ctx, `UPDATE workflows SET name=$2,description=$3,trigger=$4,trigger_cron=$5,enabled=$6,graph_draft=$7,graph_content=$8,has_draft=$9,has_content=$10,updated_at=$11 WHERE id=$1`, in.ID, in.Name, in.Description, in.Trigger, in.TriggerCron, in.Enabled, draftRaw, contentRaw, in.HasDraft, in.HasContent, now)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *Repository) DeleteWorkflow(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM workflows WHERE id=$1`, id)
	return err
}

func (r *Repository) CreateWorkflowRun(ctx context.Context, run model.WorkflowRun) (*model.WorkflowRun, error) {
	now := time.Now()
	raw, _ := json.Marshal(run.Graph)
	run.ID = uuid.NewString()
	run.CreatedAt = now
	run.UpdatedAt = now
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO workflow_runs(id,workflow_id,status,trigger,started_at,ended_at,graph,error,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, run.ID, run.WorkflowID, run.Status, run.Trigger, run.StartedAt, run.EndedAt, raw, run.Error, now, now)
	if err != nil {
		return nil, err
	}
	_, err = r.db.Pool.Exec(ctx, `UPDATE workflows SET last_run_id=$2,last_run_status=$3,last_run_time=$4,updated_at=$4 WHERE id=$1`, run.WorkflowID, run.ID, run.Status, now)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *Repository) GetWorkflowRun(ctx context.Context, id string) (*model.WorkflowRun, error) {
	m := &model.WorkflowRun{}
	var raw []byte
	err := r.db.Pool.QueryRow(ctx, `SELECT id,workflow_id,status,trigger,started_at,ended_at,graph,error,created_at,updated_at FROM workflow_runs WHERE id=$1`, id).
		Scan(&m.ID, &m.WorkflowID, &m.Status, &m.Trigger, &m.StartedAt, &m.EndedAt, &raw, &m.Error, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &m.Graph)
	}
	return m, nil
}

func (r *Repository) UpdateWorkflowRunStatus(ctx context.Context, id, status, errMsg string) error {
	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, `UPDATE workflow_runs SET status=$2,error=$3,ended_at=$4,updated_at=$4 WHERE id=$1`, id, status, errMsg, now)
	if err != nil {
		return err
	}
	_, _ = r.db.Pool.Exec(ctx, `UPDATE workflows SET last_run_status=$2,last_run_time=$3,updated_at=$3 WHERE last_run_id=$1`, id, status, now)
	return nil
}

func (r *Repository) ListWorkflowRunsByWorkflow(ctx context.Context, workflowID string, limit int) ([]model.WorkflowRun, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id,workflow_id,status,trigger,started_at,ended_at,graph,error,created_at,updated_at FROM workflow_runs WHERE workflow_id=$1 ORDER BY created_at DESC LIMIT $2`, workflowID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.WorkflowRun, 0)
	for rows.Next() {
		var m model.WorkflowRun
		var raw []byte
		if err := rows.Scan(&m.ID, &m.WorkflowID, &m.Status, &m.Trigger, &m.StartedAt, &m.EndedAt, &raw, &m.Error, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &m.Graph)
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) ListCertificates(ctx context.Context, limit int) ([]model.Certificate, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id,source,subject_alt_names,serial_number,certificate,private_key,issuer_org,key_algorithm,validity_not_after,is_revoked,workflow_id,workflow_run_id,created_at,updated_at FROM certificates ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Certificate, 0)
	for rows.Next() {
		var m model.Certificate
		if err := rows.Scan(&m.ID, &m.Source, &m.SubjectAltNames, &m.SerialNumber, &m.Certificate, &m.PrivateKey, &m.IssuerOrg, &m.KeyAlgorithm, &m.ValidityNotAfter, &m.IsRevoked, &m.WorkflowID, &m.WorkflowRunID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) GetCertificate(ctx context.Context, id string) (*model.Certificate, error) {
	m := &model.Certificate{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id,source,subject_alt_names,serial_number,certificate,private_key,issuer_org,key_algorithm,validity_not_after,is_revoked,workflow_id,workflow_run_id,created_at,updated_at FROM certificates WHERE id=$1`, id).
		Scan(&m.ID, &m.Source, &m.SubjectAltNames, &m.SerialNumber, &m.Certificate, &m.PrivateKey, &m.IssuerOrg, &m.KeyAlgorithm, &m.ValidityNotAfter, &m.IsRevoked, &m.WorkflowID, &m.WorkflowRunID, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *Repository) SaveCertificate(ctx context.Context, in model.Certificate) (*model.Certificate, error) {
	now := time.Now()
	if in.ID == "" {
		in.ID = uuid.NewString()
		in.CreatedAt = now
		in.UpdatedAt = now
		_, err := r.db.Pool.Exec(ctx, `INSERT INTO certificates(id,source,subject_alt_names,serial_number,certificate,private_key,issuer_org,key_algorithm,validity_not_after,is_revoked,workflow_id,workflow_run_id,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`, in.ID, in.Source, in.SubjectAltNames, in.SerialNumber, in.Certificate, in.PrivateKey, in.IssuerOrg, in.KeyAlgorithm, in.ValidityNotAfter, in.IsRevoked, in.WorkflowID, in.WorkflowRunID, now, now)
		if err != nil {
			return nil, err
		}
		return &in, nil
	}
	in.UpdatedAt = now
	_, err := r.db.Pool.Exec(ctx, `UPDATE certificates SET source=$2,subject_alt_names=$3,serial_number=$4,certificate=$5,private_key=$6,issuer_org=$7,key_algorithm=$8,validity_not_after=$9,is_revoked=$10,workflow_id=$11,workflow_run_id=$12,updated_at=$13 WHERE id=$1`, in.ID, in.Source, in.SubjectAltNames, in.SerialNumber, in.Certificate, in.PrivateKey, in.IssuerOrg, in.KeyAlgorithm, in.ValidityNotAfter, in.IsRevoked, in.WorkflowID, in.WorkflowRunID, now)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *Repository) RevokeCertificate(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE certificates SET is_revoked=true,updated_at=$2 WHERE id=$1`, id, time.Now())
	return err
}

func (r *Repository) GetStatistics(ctx context.Context) (*model.Statistics, error) {
	st := &model.Statistics{}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM certificates`).Scan(&st.CertificateTotal); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM certificates WHERE validity_not_after IS NOT NULL AND validity_not_after < now() + interval '21 day' AND validity_not_after >= now()`).Scan(&st.CertificateExpiringSoon); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM certificates WHERE validity_not_after IS NOT NULL AND validity_not_after < now()`).Scan(&st.CertificateExpired); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflows`).Scan(&st.WorkflowTotal); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflows WHERE enabled=true`).Scan(&st.WorkflowEnabled); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflows WHERE enabled=false`).Scan(&st.WorkflowDisabled); err != nil {
		return nil, err
	}
	return st, nil
}
