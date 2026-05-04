package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

func (r *Repository) EnsureWorkflowRunTables(ctx context.Context) error {
	required := []string{"workflow_run_nodes", "workflow_run_events", "api_keys"}
	for _, table := range required {
		var regClass *string
		if err := r.db.Pool.QueryRow(ctx, `SELECT to_regclass($1)`, "public."+table).Scan(&regClass); err != nil {
			return err
		}
		if regClass == nil || *regClass == "" {
			return fmt.Errorf("missing required table %s: run `cd server && go run ./cmd/migrate`", table)
		}
	}
	return nil
}

func isAdmin(role string) bool {
	return strings.EqualFold(strings.TrimSpace(role), model.RoleAdmin)
}

func scopedWhere(base string, userID, role string, argPos int) (string, []interface{}, int) {
	if isAdmin(role) {
		return base, nil, argPos
	}
	if strings.TrimSpace(userID) == "" {
		return base + " AND owner_user_id=''", nil, argPos
	}
	return base + fmt.Sprintf(" AND owner_user_id=$%d", argPos), []interface{}{userID}, argPos + 1
}

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash, role string) (*model.User, error) {
	if role == "" {
		role = model.RoleAdmin
	}
	id := uuid.NewString()
	now := time.Now()
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO admins(id,email,password_hash,role,status,created_at,updated_at) VALUES($1,$2,$3,$4,'active',$5,$6)`, id, email, passwordHash, role, now, now)
	if err != nil {
		return nil, err
	}
	return &model.User{ID: id, Email: email, Role: role, Status: "active", PasswordHash: passwordHash, CreatedAt: now, UpdatedAt: now}, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	m := &model.User{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id,email,role,status,password_hash,created_at,updated_at FROM admins WHERE email=$1`, email).
		Scan(&m.ID, &m.Email, &m.Role, &m.Status, &m.PasswordHash, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	m := &model.User{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id,email,role,status,password_hash,created_at,updated_at FROM admins WHERE id=$1`, id).
		Scan(&m.ID, &m.Email, &m.Role, &m.Status, &m.PasswordHash, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *Repository) UpdateUserPassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE admins SET password_hash=$2,updated_at=$3 WHERE id=$1`, id, passwordHash, time.Now())
	return err
}

func (r *Repository) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id,email,role,status,password_hash,created_at,updated_at FROM admins ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.User, 0)
	for rows.Next() {
		var m model.User
		if err := rows.Scan(&m.ID, &m.Email, &m.Role, &m.Status, &m.PasswordHash, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) UpdateUserStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE admins SET status=$2,updated_at=$3 WHERE id=$1`, id, status, time.Now())
	return err
}

func (r *Repository) ListAccessesForUser(ctx context.Context, userID, role string) ([]model.Access, error) {
	query, extraArgs, _ := scopedWhere(`SELECT id,owner_user_id,name,provider,config,reserve,deleted_at,created_at,updated_at FROM accesses WHERE deleted_at IS NULL`, userID, role, 1)
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, extraArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.Access, 0)
	for rows.Next() {
		var m model.Access
		var raw []byte
		if err := rows.Scan(&m.ID, &m.OwnerUserID, &m.Name, &m.Provider, &raw, &m.Reserve, &m.DeletedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &m.Config)
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) GetAccessByID(ctx context.Context, id string) (*model.Access, error) {
	m := &model.Access{}
	var raw []byte
	err := r.db.Pool.QueryRow(ctx, `SELECT id,owner_user_id,name,provider,config,reserve,deleted_at,created_at,updated_at FROM accesses WHERE id=$1 AND deleted_at IS NULL`, id).
		Scan(&m.ID, &m.OwnerUserID, &m.Name, &m.Provider, &raw, &m.Reserve, &m.DeletedAt, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &m.Config)
	}
	return m, nil
}

func (r *Repository) GetAccessByIDForUser(ctx context.Context, id, userID, role string) (*model.Access, error) {
	m, err := r.GetAccessByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if isAdmin(role) || m.OwnerUserID == userID {
		return m, nil
	}
	return nil, ErrNotFound
}

func (r *Repository) SaveAccessForUser(ctx context.Context, in model.Access, userID, role string) (*model.Access, error) {
	now := time.Now()
	raw, _ := json.Marshal(in.Config)
	if in.ID == "" {
		in.ID = uuid.NewString()
		in.OwnerUserID = userID
		in.CreatedAt = now
		in.UpdatedAt = now
		_, err := r.db.Pool.Exec(ctx, `INSERT INTO accesses(id,owner_user_id,name,provider,config,reserve,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, in.ID, in.OwnerUserID, in.Name, in.Provider, raw, in.Reserve, now, now)
		if err != nil {
			return nil, err
		}
		return &in, nil
	}
	if _, err := r.GetAccessByIDForUser(ctx, in.ID, userID, role); err != nil {
		return nil, err
	}
	in.UpdatedAt = now
	_, err := r.db.Pool.Exec(ctx, `UPDATE accesses SET name=$2,provider=$3,config=$4,reserve=$5,updated_at=$6 WHERE id=$1`, in.ID, in.Name, in.Provider, raw, in.Reserve, now)
	if err != nil {
		return nil, err
	}
	saved, err := r.GetAccessByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (r *Repository) SoftDeleteAccessForUser(ctx context.Context, id, userID, role string) error {
	if _, err := r.GetAccessByIDForUser(ctx, id, userID, role); err != nil {
		return err
	}
	_, err := r.db.Pool.Exec(ctx, `UPDATE accesses SET deleted_at=$2,updated_at=$2 WHERE id=$1`, id, time.Now())
	return err
}

func (r *Repository) ListWorkflowsForUser(ctx context.Context, userID, role string) ([]model.Workflow, error) {
	query, extraArgs, _ := scopedWhere(`SELECT id,owner_user_id,name,description,trigger,trigger_cron,enabled,graph_draft,graph_content,has_draft,has_content,last_run_id,last_run_status,last_run_time,created_at,updated_at FROM workflows WHERE 1=1`, userID, role, 1)
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, extraArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Workflow, 0)
	for rows.Next() {
		var m model.Workflow
		var draftRaw, contentRaw []byte
		if err := rows.Scan(&m.ID, &m.OwnerUserID, &m.Name, &m.Description, &m.Trigger, &m.TriggerCron, &m.Enabled, &draftRaw, &contentRaw, &m.HasDraft, &m.HasContent, &m.LastRunID, &m.LastRunStatus, &m.LastRunTime, &m.CreatedAt, &m.UpdatedAt); err != nil {
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
	err := r.db.Pool.QueryRow(ctx, `SELECT id,owner_user_id,name,description,trigger,trigger_cron,enabled,graph_draft,graph_content,has_draft,has_content,last_run_id,last_run_status,last_run_time,created_at,updated_at FROM workflows WHERE id=$1`, id).
		Scan(&m.ID, &m.OwnerUserID, &m.Name, &m.Description, &m.Trigger, &m.TriggerCron, &m.Enabled, &draftRaw, &contentRaw, &m.HasDraft, &m.HasContent, &m.LastRunID, &m.LastRunStatus, &m.LastRunTime, &m.CreatedAt, &m.UpdatedAt)
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

func (r *Repository) GetWorkflowForUser(ctx context.Context, id, userID, role string) (*model.Workflow, error) {
	m, err := r.GetWorkflow(ctx, id)
	if err != nil {
		return nil, err
	}
	if isAdmin(role) || m.OwnerUserID == userID {
		return m, nil
	}
	return nil, ErrNotFound
}

func (r *Repository) SaveWorkflowForUser(ctx context.Context, in model.Workflow, userID, role string) (*model.Workflow, error) {
	now := time.Now()
	draftRaw, _ := json.Marshal(in.GraphDraft)
	contentRaw, _ := json.Marshal(in.GraphContent)
	if in.ID == "" {
		in.ID = uuid.NewString()
		in.OwnerUserID = userID
		in.CreatedAt = now
		in.UpdatedAt = now
		_, err := r.db.Pool.Exec(ctx, `INSERT INTO workflows(id,owner_user_id,name,description,trigger,trigger_cron,enabled,graph_draft,graph_content,has_draft,has_content,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`, in.ID, in.OwnerUserID, in.Name, in.Description, in.Trigger, in.TriggerCron, in.Enabled, draftRaw, contentRaw, in.HasDraft, in.HasContent, now, now)
		if err != nil {
			return nil, err
		}
		return &in, nil
	}
	if _, err := r.GetWorkflowForUser(ctx, in.ID, userID, role); err != nil {
		return nil, err
	}
	in.UpdatedAt = now
	_, err := r.db.Pool.Exec(ctx, `UPDATE workflows SET name=$2,description=$3,trigger=$4,trigger_cron=$5,enabled=$6,graph_draft=$7,graph_content=$8,has_draft=$9,has_content=$10,updated_at=$11 WHERE id=$1`, in.ID, in.Name, in.Description, in.Trigger, in.TriggerCron, in.Enabled, draftRaw, contentRaw, in.HasDraft, in.HasContent, now)
	if err != nil {
		return nil, err
	}
	return r.GetWorkflow(ctx, in.ID)
}

func (r *Repository) DeleteWorkflowForUser(ctx context.Context, id, userID, role string) error {
	if _, err := r.GetWorkflowForUser(ctx, id, userID, role); err != nil {
		return err
	}
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM workflows WHERE id=$1`, id)
	return err
}

func (r *Repository) CreateWorkflowRun(ctx context.Context, run model.WorkflowRun) (*model.WorkflowRun, error) {
	now := time.Now()
	raw, _ := json.Marshal(run.Graph)
	run.ID = uuid.NewString()
	run.CreatedAt = now
	run.UpdatedAt = now
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO workflow_runs(id,workflow_id,owner_user_id,status,trigger,started_at,ended_at,graph,error,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, run.ID, run.WorkflowID, run.OwnerUserID, run.Status, run.Trigger, run.StartedAt, run.EndedAt, raw, run.Error, now, now)
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
	err := r.db.Pool.QueryRow(ctx, `SELECT id,workflow_id,owner_user_id,status,trigger,started_at,ended_at,graph,error,created_at,updated_at FROM workflow_runs WHERE id=$1`, id).
		Scan(&m.ID, &m.WorkflowID, &m.OwnerUserID, &m.Status, &m.Trigger, &m.StartedAt, &m.EndedAt, &raw, &m.Error, &m.CreatedAt, &m.UpdatedAt)
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

func (r *Repository) ListWorkflowRunsByWorkflowForUser(ctx context.Context, workflowID, userID, role string, limit int) ([]model.WorkflowRun, error) {
	if _, err := r.GetWorkflowForUser(ctx, workflowID, userID, role); err != nil {
		return nil, err
	}

	query, extraArgs, argPos := scopedWhere(`SELECT id,workflow_id,owner_user_id,status,trigger,started_at,ended_at,graph,error,created_at,updated_at FROM workflow_runs WHERE workflow_id=$1`, userID, role, 2)
	args := []interface{}{workflowID}
	args = append(args, extraArgs...)
	args = append(args, limit)
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argPos)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.WorkflowRun, 0)
	for rows.Next() {
		var m model.WorkflowRun
		var raw []byte
		if err := rows.Scan(&m.ID, &m.WorkflowID, &m.OwnerUserID, &m.Status, &m.Trigger, &m.StartedAt, &m.EndedAt, &raw, &m.Error, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			_ = json.Unmarshal(raw, &m.Graph)
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) GetWorkflowRunForUser(ctx context.Context, runID, userID, role string) (*model.WorkflowRun, error) {
	run, err := r.GetWorkflowRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	if isAdmin(role) || run.OwnerUserID == userID {
		return run, nil
	}
	return nil, ErrNotFound
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

func (r *Repository) UpsertWorkflowRunNode(ctx context.Context, in model.WorkflowRunNode) (*model.WorkflowRunNode, error) {
	now := time.Now()
	outputRaw, _ := json.Marshal(in.Output)

	if in.ID == "" {
		in.ID = uuid.NewString()
	}
	in.UpdatedAt = now

	_, err := r.db.Pool.Exec(ctx, `
INSERT INTO workflow_run_nodes(id,run_id,node_id,node_name,action,provider,status,started_at,ended_at,error,output,created_at,updated_at)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
ON CONFLICT(run_id,node_id) DO UPDATE SET
  node_name=EXCLUDED.node_name,
  action=EXCLUDED.action,
  provider=EXCLUDED.provider,
  status=EXCLUDED.status,
  started_at=COALESCE(EXCLUDED.started_at, workflow_run_nodes.started_at),
  ended_at=EXCLUDED.ended_at,
  error=EXCLUDED.error,
  output=EXCLUDED.output,
  updated_at=EXCLUDED.updated_at
`, in.ID, in.RunID, in.NodeID, in.NodeName, in.Action, in.Provider, in.Status, in.StartedAt, in.EndedAt, in.Error, outputRaw, now, now)
	if err != nil {
		return nil, err
	}
	return r.GetWorkflowRunNode(ctx, in.RunID, in.NodeID)
}

func (r *Repository) GetWorkflowRunNode(ctx context.Context, runID, nodeID string) (*model.WorkflowRunNode, error) {
	var m model.WorkflowRunNode
	var outputRaw []byte
	err := r.db.Pool.QueryRow(ctx, `SELECT id,run_id,node_id,node_name,action,provider,status,started_at,ended_at,error,output,created_at,updated_at FROM workflow_run_nodes WHERE run_id=$1 AND node_id=$2`, runID, nodeID).
		Scan(&m.ID, &m.RunID, &m.NodeID, &m.NodeName, &m.Action, &m.Provider, &m.Status, &m.StartedAt, &m.EndedAt, &m.Error, &outputRaw, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if len(outputRaw) > 0 {
		_ = json.Unmarshal(outputRaw, &m.Output)
	}
	return &m, nil
}

func (r *Repository) ListWorkflowRunNodesForUser(ctx context.Context, runID, userID, role string) ([]model.WorkflowRunNode, error) {
	if _, err := r.GetWorkflowRunForUser(ctx, runID, userID, role); err != nil {
		return nil, err
	}
	rows, err := r.db.Pool.Query(ctx, `SELECT id,run_id,node_id,node_name,action,provider,status,started_at,ended_at,error,output,created_at,updated_at FROM workflow_run_nodes WHERE run_id=$1 ORDER BY created_at ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.WorkflowRunNode, 0)
	for rows.Next() {
		var m model.WorkflowRunNode
		var outputRaw []byte
		if err := rows.Scan(&m.ID, &m.RunID, &m.NodeID, &m.NodeName, &m.Action, &m.Provider, &m.Status, &m.StartedAt, &m.EndedAt, &m.Error, &outputRaw, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		if len(outputRaw) > 0 {
			_ = json.Unmarshal(outputRaw, &m.Output)
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) AppendWorkflowRunEvent(ctx context.Context, in model.WorkflowRunEvent) (*model.WorkflowRunEvent, error) {
	if in.ID == "" {
		in.ID = uuid.NewString()
	}
	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}
	payloadRaw, _ := json.Marshal(in.Payload)
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO workflow_run_events(id,run_id,node_id,event_type,message,payload,created_at) VALUES($1,$2,$3,$4,$5,$6,$7)`, in.ID, in.RunID, in.NodeID, in.EventType, in.Message, payloadRaw, in.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *Repository) ListWorkflowRunEventsForUser(ctx context.Context, runID, userID, role, nodeID string, since *time.Time, limit int) ([]model.WorkflowRunEvent, error) {
	if _, err := r.GetWorkflowRunForUser(ctx, runID, userID, role); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	where := "WHERE run_id=$1"
	args := []interface{}{runID}
	argPos := 2
	if nodeID != "" {
		where += fmt.Sprintf(" AND node_id=$%d", argPos)
		args = append(args, nodeID)
		argPos++
	}
	if since != nil {
		where += fmt.Sprintf(" AND created_at>$%d", argPos)
		args = append(args, *since)
		argPos++
	}
	args = append(args, limit)
	query := `SELECT id,run_id,node_id,event_type,message,payload,created_at FROM workflow_run_events ` + where + ` ORDER BY created_at ASC LIMIT $` + fmt.Sprintf("%d", argPos)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.WorkflowRunEvent, 0)
	for rows.Next() {
		var m model.WorkflowRunEvent
		var payloadRaw []byte
		if err := rows.Scan(&m.ID, &m.RunID, &m.NodeID, &m.EventType, &m.Message, &payloadRaw, &m.CreatedAt); err != nil {
			return nil, err
		}
		if len(payloadRaw) > 0 {
			_ = json.Unmarshal(payloadRaw, &m.Payload)
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) ListCertificatesForUser(ctx context.Context, userID, role string, limit int) ([]model.Certificate, error) {
	query, extraArgs, argPos := scopedWhere(`SELECT id,owner_user_id,source,subject_alt_names,serial_number,certificate,private_key,issuer_org,key_algorithm,validity_not_after,is_revoked,workflow_id,workflow_run_id,created_at,updated_at FROM certificates WHERE 1=1`, userID, role, 1)
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d`, argPos)
	args := append(extraArgs, limit)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Certificate, 0)
	for rows.Next() {
		var m model.Certificate
		if err := rows.Scan(&m.ID, &m.OwnerUserID, &m.Source, &m.SubjectAltNames, &m.SerialNumber, &m.Certificate, &m.PrivateKey, &m.IssuerOrg, &m.KeyAlgorithm, &m.ValidityNotAfter, &m.IsRevoked, &m.WorkflowID, &m.WorkflowRunID, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) GetCertificate(ctx context.Context, id string) (*model.Certificate, error) {
	m := &model.Certificate{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id,owner_user_id,source,subject_alt_names,serial_number,certificate,private_key,issuer_org,key_algorithm,validity_not_after,is_revoked,workflow_id,workflow_run_id,created_at,updated_at FROM certificates WHERE id=$1`, id).
		Scan(&m.ID, &m.OwnerUserID, &m.Source, &m.SubjectAltNames, &m.SerialNumber, &m.Certificate, &m.PrivateKey, &m.IssuerOrg, &m.KeyAlgorithm, &m.ValidityNotAfter, &m.IsRevoked, &m.WorkflowID, &m.WorkflowRunID, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *Repository) GetCertificateForUser(ctx context.Context, id, userID, role string) (*model.Certificate, error) {
	m, err := r.GetCertificate(ctx, id)
	if err != nil {
		return nil, err
	}
	if isAdmin(role) || m.OwnerUserID == userID {
		return m, nil
	}
	return nil, ErrNotFound
}

func (r *Repository) GetLatestCertificateByRunForUser(ctx context.Context, runID, userID, role string) (*model.Certificate, error) {
	query, extraArgs, _ := scopedWhere(`SELECT id,owner_user_id,source,subject_alt_names,serial_number,certificate,private_key,issuer_org,key_algorithm,validity_not_after,is_revoked,workflow_id,workflow_run_id,created_at,updated_at FROM certificates WHERE workflow_run_id=$1`, userID, role, 2)
	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT 1`)
	args := []interface{}{runID}
	args = append(args, extraArgs...)

	m := &model.Certificate{}
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&m.ID, &m.OwnerUserID, &m.Source, &m.SubjectAltNames, &m.SerialNumber, &m.Certificate, &m.PrivateKey, &m.IssuerOrg, &m.KeyAlgorithm, &m.ValidityNotAfter, &m.IsRevoked, &m.WorkflowID, &m.WorkflowRunID, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *Repository) SaveCertificate(ctx context.Context, in model.Certificate) (*model.Certificate, error) {
	now := time.Now()
	if in.ID == "" {
		in.ID = uuid.NewString()
		in.CreatedAt = now
		in.UpdatedAt = now
		_, err := r.db.Pool.Exec(ctx, `INSERT INTO certificates(id,owner_user_id,source,subject_alt_names,serial_number,certificate,private_key,issuer_org,key_algorithm,validity_not_after,is_revoked,workflow_id,workflow_run_id,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`, in.ID, in.OwnerUserID, in.Source, in.SubjectAltNames, in.SerialNumber, in.Certificate, in.PrivateKey, in.IssuerOrg, in.KeyAlgorithm, in.ValidityNotAfter, in.IsRevoked, in.WorkflowID, in.WorkflowRunID, now, now)
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

func (r *Repository) RevokeCertificateForUser(ctx context.Context, id, userID, role string) error {
	if _, err := r.GetCertificateForUser(ctx, id, userID, role); err != nil {
		return err
	}
	_, err := r.db.Pool.Exec(ctx, `UPDATE certificates SET is_revoked=true,updated_at=$2 WHERE id=$1`, id, time.Now())
	return err
}

func (r *Repository) GetStatisticsForUser(ctx context.Context, userID, role string) (*model.Statistics, error) {
	st := &model.Statistics{}

	certWhere := ""
	wfWhere := ""
	args := []interface{}{}
	if !isAdmin(role) {
		certWhere = " WHERE owner_user_id=$1"
		wfWhere = " WHERE owner_user_id=$1"
		args = append(args, userID)
	}

	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM certificates`+certWhere, args...).Scan(&st.CertificateTotal); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM certificates`+certWhere+andClause(certWhere)+` validity_not_after IS NOT NULL AND validity_not_after < now() + interval '21 day' AND validity_not_after >= now()`, args...).Scan(&st.CertificateExpiringSoon); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM certificates`+certWhere+andClause(certWhere)+` validity_not_after IS NOT NULL AND validity_not_after < now()`, args...).Scan(&st.CertificateExpired); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflows`+wfWhere, args...).Scan(&st.WorkflowTotal); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflows`+wfWhere+andClause(wfWhere)+` enabled=true`, args...).Scan(&st.WorkflowEnabled); err != nil {
		return nil, err
	}
	if err := r.db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflows`+wfWhere+andClause(wfWhere)+` enabled=false`, args...).Scan(&st.WorkflowDisabled); err != nil {
		return nil, err
	}
	return st, nil
}

func andClause(where string) string {
	if where == "" {
		return " WHERE"
	}
	return " AND"
}

func (r *Repository) CreateAPIKey(ctx context.Context, in model.APIKey) (*model.APIKey, error) {
	now := time.Now()
	if in.ID == "" {
		in.ID = uuid.NewString()
	}
	in.CreatedAt = now
	in.UpdatedAt = now
	if in.Status == "" {
		in.Status = "active"
	}
	_, err := r.db.Pool.Exec(ctx, `INSERT INTO api_keys(id,user_id,name,prefix,key_hash,status,expires_at,last_used_at,revoked_at,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, in.ID, in.UserID, in.Name, in.Prefix, in.KeyHash, in.Status, in.ExpiresAt, in.LastUsedAt, in.RevokedAt, in.CreatedAt, in.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *Repository) ListAPIKeysByUser(ctx context.Context, userID string) ([]model.APIKey, error) {
	rows, err := r.db.Pool.Query(ctx, `SELECT id,user_id,name,prefix,key_hash,status,expires_at,last_used_at,revoked_at,created_at,updated_at FROM api_keys WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.APIKey, 0)
	for rows.Next() {
		var m model.APIKey
		if err := rows.Scan(&m.ID, &m.UserID, &m.Name, &m.Prefix, &m.KeyHash, &m.Status, &m.ExpiresAt, &m.LastUsedAt, &m.RevokedAt, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *Repository) GetAPIKeyByPrefix(ctx context.Context, prefix string) (*model.APIKey, error) {
	m := &model.APIKey{}
	err := r.db.Pool.QueryRow(ctx, `SELECT id,user_id,name,prefix,key_hash,status,expires_at,last_used_at,revoked_at,created_at,updated_at FROM api_keys WHERE prefix=$1`, prefix).
		Scan(&m.ID, &m.UserID, &m.Name, &m.Prefix, &m.KeyHash, &m.Status, &m.ExpiresAt, &m.LastUsedAt, &m.RevokedAt, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return m, err
}

func (r *Repository) RevokeAPIKey(ctx context.Context, id, userID string) error {
	res, err := r.db.Pool.Exec(ctx, `UPDATE api_keys SET status='revoked',revoked_at=$3,updated_at=$3 WHERE id=$1 AND user_id=$2 AND status='active'`, id, userID, time.Now())
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) TouchAPIKeyUsage(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, `UPDATE api_keys SET last_used_at=$2,updated_at=$2 WHERE id=$1`, id, time.Now())
	return err
}
