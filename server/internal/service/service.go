package service

import (
	"context"
	"errors"
	"time"

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
	return s.repo.ListAccesses(ctx)
}
func (s *Service) SaveAccess(ctx context.Context, in model.Access) (*model.Access, error) {
	return s.repo.SaveAccess(ctx, in)
}
func (s *Service) DeleteAccess(ctx context.Context, id string) error {
	return s.repo.SoftDeleteAccess(ctx, id)
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
	file := c.Certificate + "\n" + c.PrivateKey
	if format == "JKS" || format == "PFX" {
		file = "MOCK-" + format + "-BINARY"
	}
	return map[string]interface{}{"fileBytes": file, "fileFormat": format}, nil
}

func (s *Service) Statistics(ctx context.Context) (*model.Statistics, error) {
	return s.repo.GetStatistics(ctx)
}

func (s *Service) TestNotification(ctx context.Context, provider, accessID string) error {
	_ = provider
	_ = accessID
	return nil
}
