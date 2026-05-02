package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"easyssl/server/internal/accessprovider"
	"easyssl/server/internal/middleware"
	"easyssl/server/internal/model"
	"easyssl/server/internal/repository"
	wf "easyssl/server/internal/workflow"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo       *repository.Repository
	jwtSecret  string
	dispatcher *wf.Dispatcher
}

func New(repo *repository.Repository, jwtSecret string, dispatcher *wf.Dispatcher) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret, dispatcher: dispatcher}
}

func (s *Service) EnsureBootstrapAdmin(ctx context.Context, email, password string) error {
	_, err := s.repo.GetAdminByEmail(ctx, email)
	if err == nil {
		return nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.repo.CreateAdmin(ctx, email, string(h))
	return err
}

func (s *Service) Login(ctx context.Context, email, password string) (string, *model.Admin, error) {
	admin, err := s.repo.GetAdminByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	tk, err := middleware.Sign(s.jwtSecret, admin.ID)
	if err != nil {
		return "", nil, err
	}
	return tk, admin, nil
}

func (s *Service) Me(ctx context.Context, adminID string) (*model.Admin, error) {
	return s.repo.GetAdminByID(ctx, adminID)
}

func (s *Service) ChangePassword(ctx context.Context, adminID, password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdateAdminPassword(ctx, adminID, string(h))
}

func (s *Service) ListAccesses(ctx context.Context) ([]model.Access, error) {
	items, err := s.repo.ListAccesses(ctx)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i] = sanitizeAccess(items[i])
	}
	return items, nil
}
func (s *Service) SaveAccess(ctx context.Context, in model.Access) (*model.Access, error) {
	normalized, err := s.prepareAccess(ctx, in)
	if err != nil {
		return nil, err
	}
	saved, err := s.repo.SaveAccess(ctx, normalized)
	if err != nil {
		return nil, err
	}
	out := sanitizeAccess(*saved)
	return &out, nil
}
func (s *Service) DeleteAccess(ctx context.Context, id string) error {
	return s.repo.SoftDeleteAccess(ctx, id)
}
func (s *Service) TestAccess(ctx context.Context, id string) error {
	access, err := s.repo.GetAccessByID(ctx, id)
	if err != nil {
		return err
	}
	return accessprovider.TestAccess(ctx, *access)
}

func (s *Service) ListWorkflows(ctx context.Context) ([]model.Workflow, error) {
	return s.repo.ListWorkflows(ctx)
}
func (s *Service) GetWorkflow(ctx context.Context, id string) (*model.Workflow, error) {
	return s.repo.GetWorkflow(ctx, id)
}
func (s *Service) SaveWorkflow(ctx context.Context, in model.Workflow) (*model.Workflow, error) {
	return s.repo.SaveWorkflow(ctx, in)
}
func (s *Service) DeleteWorkflow(ctx context.Context, id string) error {
	return s.repo.DeleteWorkflow(ctx, id)
}
func (s *Service) ListWorkflowRuns(ctx context.Context, workflowID string) ([]model.WorkflowRun, error) {
	return s.repo.ListWorkflowRunsByWorkflow(ctx, workflowID, 30)
}

func (s *Service) ListWorkflowRunNodes(ctx context.Context, workflowID, runID string) ([]model.WorkflowRunNode, error) {
	run, err := s.repo.GetWorkflowRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run.WorkflowID != workflowID {
		return nil, fmt.Errorf("run %s not found in workflow %s", runID, workflowID)
	}
	return s.repo.ListWorkflowRunNodes(ctx, runID)
}

func (s *Service) ListWorkflowRunEvents(ctx context.Context, workflowID, runID, nodeID string, since *time.Time, limit int) ([]model.WorkflowRunEvent, error) {
	run, err := s.repo.GetWorkflowRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run.WorkflowID != workflowID {
		return nil, fmt.Errorf("run %s not found in workflow %s", runID, workflowID)
	}
	return s.repo.ListWorkflowRunEvents(ctx, runID, nodeID, since, limit)
}

func (s *Service) StartWorkflowRun(ctx context.Context, workflowID, trigger string) (string, error) {
	w, err := s.repo.GetWorkflow(ctx, workflowID)
	if err != nil {
		return "", err
	}
	run, err := s.repo.CreateWorkflowRun(ctx, model.WorkflowRun{WorkflowID: workflowID, Status: "pending", Trigger: trigger, StartedAt: time.Now(), Graph: w.GraphContent})
	if err != nil {
		return "", err
	}
	s.dispatcher.Start(ctx, run.ID)
	return run.ID, nil
}

func (s *Service) CancelWorkflowRun(ctx context.Context, runID string) error {
	s.dispatcher.Cancel(ctx, runID)
	return nil
}

func (s *Service) WorkflowStats() map[string]interface{} {
	c, pending, processing := s.dispatcher.Stats()
	return map[string]interface{}{"concurrency": c, "pendingRunIds": pending, "processingRunIds": processing}
}

func (s *Service) ListCertificates(ctx context.Context) ([]model.Certificate, error) {
	return s.repo.ListCertificates(ctx, 200)
}
func (s *Service) RevokeCertificate(ctx context.Context, id string) error {
	return s.repo.RevokeCertificate(ctx, id)
}

func (s *Service) DownloadCertificate(ctx context.Context, id, format string) (map[string]interface{}, error) {
	c, err := s.repo.GetCertificate(ctx, id)
	if err != nil {
		return nil, err
	}
	switch format {
	case "", "PEM":
		zipBytes, err := buildPEMArchive(c.Certificate, c.PrivateKey)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"fileName":        c.ID + ".zip",
			"fileFormat":      "ZIP",
			"mimeType":        "application/zip",
			"fileBytesBase64": base64.StdEncoding.EncodeToString(zipBytes),
		}, nil
	case "PFX", "JKS":
		content := []byte("MOCK-" + format + "-BINARY")
		return map[string]interface{}{
			"fileName":        c.ID + "." + lower(format),
			"fileFormat":      format,
			"mimeType":        "application/octet-stream",
			"fileBytesBase64": base64.StdEncoding.EncodeToString(content),
		}, nil
	default:
		return nil, errors.New("unsupported certificate format")
	}
}

func buildPEMArchive(fullchainPEM, privateKeyPEM string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	writeText := func(name, content string) error {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(content))
		return err
	}

	if err := writeText("fullchain.pem", fullchainPEM); err != nil {
		return nil, err
	}
	if err := writeText("privkey.pem", privateKeyPEM); err != nil {
		return nil, err
	}

	serverCert, chainCert := splitServerAndChain(fullchainPEM)
	if err := writeText("cert.pem", serverCert); err != nil {
		return nil, err
	}
	if chainCert != "" {
		if err := writeText("chain.pem", chainCert); err != nil {
			return nil, err
		}
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func splitServerAndChain(fullchain string) (string, string) {
	const endMark = "-----END CERTIFICATE-----"
	parts := make([]string, 0)
	remain := fullchain
	for {
		i := bytes.Index([]byte(remain), []byte(endMark))
		if i < 0 {
			break
		}
		i += len(endMark)
		block := remain[:i]
		parts = append(parts, block)
		remain = remain[i:]
	}
	if len(parts) == 0 {
		return fullchain, ""
	}
	server := parts[0]
	chain := ""
	for i := 1; i < len(parts); i++ {
		if chain == "" {
			chain = parts[i]
		} else {
			chain += "\n" + parts[i]
		}
	}
	return server, chain
}

func lower(s string) string {
	if s == "PFX" {
		return "pfx"
	}
	if s == "JKS" {
		return "jks"
	}
	return s
}

func (s *Service) Statistics(ctx context.Context) (*model.Statistics, error) {
	return s.repo.GetStatistics(ctx)
}

func (s *Service) TestNotification(ctx context.Context, provider, accessID string) error {
	_ = provider
	_ = accessID
	return nil
}
