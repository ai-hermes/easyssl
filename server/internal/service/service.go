package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"easyssl/server/internal/accessprovider"
	"easyssl/server/internal/middleware"
	"easyssl/server/internal/model"
	"easyssl/server/internal/providercatalog"
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

func (s *Service) EnsureBootstrapUser(ctx context.Context, email, password string) error {
	_, err := s.repo.GetUserByEmail(ctx, email)
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
	_, err = s.repo.CreateUser(ctx, email, string(h), model.RoleAdmin)
	return err
}

func (s *Service) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user.Status != "active" {
		return "", nil, errors.New("account is disabled")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	tk, err := middleware.Sign(s.jwtSecret, user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}
	return tk, user, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) ChangePassword(ctx context.Context, userID, password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserPassword(ctx, userID, string(h))
}

func (s *Service) Register(ctx context.Context, email, password string) (*model.User, error) {
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateUser(ctx, email, string(h), model.RoleUser)
}

func (s *Service) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *Service) UpdateUserStatus(ctx context.Context, id, status string) error {
	return s.repo.UpdateUserStatus(ctx, id, status)
}

func apiKeyHash(secret, raw string) string {
	sum := sha256.Sum256([]byte(secret + ":" + raw))
	return hex.EncodeToString(sum[:])
}

func apiKeyPrefix(raw string) string {
	if len(raw) <= 16 {
		return raw
	}
	return raw[:16]
}

func randomAPIKey() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "esk_" + base64.RawURLEncoding.EncodeToString(buf), nil
}

func (s *Service) VerifyAPIKey(ctx context.Context, rawKey string) (*model.AuthContext, error) {
	rawKey = strings.TrimSpace(rawKey)
	if rawKey == "" {
		return nil, errors.New("empty api key")
	}
	stored, err := s.repo.GetAPIKeyByPrefix(ctx, apiKeyPrefix(rawKey))
	if err != nil {
		return nil, err
	}
	if stored.Status != "active" {
		return nil, errors.New("api key revoked")
	}
	if stored.ExpiresAt != nil && stored.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("api key expired")
	}
	hashed := apiKeyHash(s.jwtSecret, rawKey)
	if subtle.ConstantTimeCompare([]byte(hashed), []byte(stored.KeyHash)) != 1 {
		return nil, errors.New("invalid api key")
	}
	user, err := s.repo.GetUserByID(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}
	if user.Status != "active" {
		return nil, errors.New("account disabled")
	}
	_ = s.repo.TouchAPIKeyUsage(ctx, stored.ID)
	return &model.AuthContext{UserID: user.ID, Role: user.Role, AuthType: "api_key", APIKeyID: stored.ID}, nil
}

func (s *Service) CreateAPIKey(ctx context.Context, auth model.AuthContext, name string, expiresAt *time.Time) (map[string]interface{}, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	raw, err := randomAPIKey()
	if err != nil {
		return nil, err
	}
	prefix := apiKeyPrefix(raw)
	saved, err := s.repo.CreateAPIKey(ctx, model.APIKey{
		UserID:    auth.UserID,
		Name:      name,
		Prefix:    prefix,
		KeyHash:   apiKeyHash(s.jwtSecret, raw),
		Status:    "active",
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":        saved.ID,
		"name":      saved.Name,
		"prefix":    saved.Prefix,
		"status":    saved.Status,
		"expiresAt": saved.ExpiresAt,
		"createdAt": saved.CreatedAt,
		"token":     raw,
	}, nil
}

func (s *Service) ListAPIKeys(ctx context.Context, auth model.AuthContext) ([]map[string]interface{}, error) {
	keys, err := s.repo.ListAPIKeysByUser(ctx, auth.UserID)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(keys))
	for _, k := range keys {
		items = append(items, map[string]interface{}{
			"id":         k.ID,
			"name":       k.Name,
			"prefix":     k.Prefix,
			"status":     k.Status,
			"expiresAt":  k.ExpiresAt,
			"lastUsedAt": k.LastUsedAt,
			"createdAt":  k.CreatedAt,
		})
	}
	return items, nil
}

func (s *Service) RevokeAPIKey(ctx context.Context, auth model.AuthContext, id string) error {
	return s.repo.RevokeAPIKey(ctx, id, auth.UserID)
}

func (s *Service) ListAccesses(ctx context.Context, auth model.AuthContext) ([]model.Access, error) {
	items, err := s.repo.ListAccessesForUser(ctx, auth.UserID, auth.Role)
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i] = sanitizeAccess(items[i])
	}
	return items, nil
}

func (s *Service) SaveAccess(ctx context.Context, auth model.AuthContext, in model.Access) (*model.Access, error) {
	normalized, err := s.prepareAccess(ctx, auth, in)
	if err != nil {
		return nil, err
	}
	saved, err := s.repo.SaveAccessForUser(ctx, normalized, auth.UserID, auth.Role)
	if err != nil {
		return nil, err
	}
	out := sanitizeAccess(*saved)
	return &out, nil
}

func (s *Service) DeleteAccess(ctx context.Context, auth model.AuthContext, id string) error {
	return s.repo.SoftDeleteAccessForUser(ctx, id, auth.UserID, auth.Role)
}

func (s *Service) TestAccess(ctx context.Context, auth model.AuthContext, id string) error {
	access, err := s.repo.GetAccessByIDForUser(ctx, id, auth.UserID, auth.Role)
	if err != nil {
		return err
	}
	return accessprovider.TestAccess(ctx, *access)
}

func (s *Service) ListWorkflows(ctx context.Context, auth model.AuthContext) ([]model.Workflow, error) {
	return s.repo.ListWorkflowsForUser(ctx, auth.UserID, auth.Role)
}

func (s *Service) GetWorkflow(ctx context.Context, auth model.AuthContext, id string) (*model.Workflow, error) {
	return s.repo.GetWorkflowForUser(ctx, id, auth.UserID, auth.Role)
}

func (s *Service) SaveWorkflow(ctx context.Context, auth model.AuthContext, in model.Workflow) (*model.Workflow, error) {
	return s.repo.SaveWorkflowForUser(ctx, in, auth.UserID, auth.Role)
}

func (s *Service) DeleteWorkflow(ctx context.Context, auth model.AuthContext, id string) error {
	return s.repo.DeleteWorkflowForUser(ctx, id, auth.UserID, auth.Role)
}

func (s *Service) ListWorkflowRuns(ctx context.Context, auth model.AuthContext, workflowID string) ([]model.WorkflowRun, error) {
	return s.repo.ListWorkflowRunsByWorkflowForUser(ctx, workflowID, auth.UserID, auth.Role, 30)
}

func (s *Service) ListWorkflowRunNodes(ctx context.Context, auth model.AuthContext, workflowID, runID string) ([]model.WorkflowRunNode, error) {
	run, err := s.repo.GetWorkflowRunForUser(ctx, runID, auth.UserID, auth.Role)
	if err != nil {
		return nil, err
	}
	if run.WorkflowID != workflowID {
		return nil, fmt.Errorf("run %s not found in workflow %s", runID, workflowID)
	}
	return s.repo.ListWorkflowRunNodesForUser(ctx, runID, auth.UserID, auth.Role)
}

func (s *Service) ListWorkflowRunEvents(ctx context.Context, auth model.AuthContext, workflowID, runID, nodeID string, since *time.Time, limit int) ([]model.WorkflowRunEvent, error) {
	run, err := s.repo.GetWorkflowRunForUser(ctx, runID, auth.UserID, auth.Role)
	if err != nil {
		return nil, err
	}
	if run.WorkflowID != workflowID {
		return nil, fmt.Errorf("run %s not found in workflow %s", runID, workflowID)
	}
	return s.repo.ListWorkflowRunEventsForUser(ctx, runID, auth.UserID, auth.Role, nodeID, since, limit)
}

func (s *Service) StartWorkflowRun(ctx context.Context, auth model.AuthContext, workflowID, trigger string) (string, error) {
	w, err := s.repo.GetWorkflowForUser(ctx, workflowID, auth.UserID, auth.Role)
	if err != nil {
		return "", err
	}
	run, err := s.repo.CreateWorkflowRun(ctx, model.WorkflowRun{OwnerUserID: w.OwnerUserID, WorkflowID: workflowID, Status: "pending", Trigger: trigger, StartedAt: time.Now(), Graph: w.GraphContent})
	if err != nil {
		return "", err
	}
	s.dispatcher.Start(ctx, run.ID)
	return run.ID, nil
}

func (s *Service) CancelWorkflowRun(ctx context.Context, auth model.AuthContext, runID string) error {
	if _, err := s.repo.GetWorkflowRunForUser(ctx, runID, auth.UserID, auth.Role); err != nil {
		return err
	}
	s.dispatcher.Cancel(ctx, runID)
	return nil
}

func (s *Service) WorkflowStats() map[string]interface{} {
	c, pending, processing := s.dispatcher.Stats()
	return map[string]interface{}{"concurrency": c, "pendingRunIds": pending, "processingRunIds": processing}
}

func (s *Service) ListCertificates(ctx context.Context, auth model.AuthContext) ([]model.Certificate, error) {
	return s.repo.ListCertificatesForUser(ctx, auth.UserID, auth.Role, 200)
}

func (s *Service) RevokeCertificate(ctx context.Context, auth model.AuthContext, id string) error {
	return s.repo.RevokeCertificateForUser(ctx, id, auth.UserID, auth.Role)
}

func (s *Service) DownloadCertificate(ctx context.Context, auth model.AuthContext, id, format string) (map[string]interface{}, error) {
	c, err := s.repo.GetCertificateForUser(ctx, id, auth.UserID, auth.Role)
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

func (s *Service) Statistics(ctx context.Context, auth model.AuthContext) (*model.Statistics, error) {
	return s.repo.GetStatisticsForUser(ctx, auth.UserID, auth.Role)
}

func (s *Service) TestNotification(ctx context.Context, auth model.AuthContext, provider, accessID string) error {
	_ = provider
	_, err := s.repo.GetAccessByIDForUser(ctx, accessID, auth.UserID, auth.Role)
	return err
}

func normalizeOpenProvider(provider string) string {
	return providercatalog.Normalize(provider)
}

func normalizeDomains(domains []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(domains))
	for _, item := range domains {
		domain := strings.TrimSpace(item)
		if domain == "" {
			continue
		}
		if _, ok := seen[domain]; ok {
			continue
		}
		seen[domain] = struct{}{}
		out = append(out, domain)
	}
	return out
}

func (s *Service) OpenApplyCertificate(ctx context.Context, auth model.AuthContext, req model.OpenApplyCertificateRequest) (map[string]interface{}, error) {
	accessID := strings.TrimSpace(req.AccessID)
	if accessID == "" {
		return nil, errors.New("accessId is required")
	}
	domains := normalizeDomains(req.Domains)
	if len(domains) == 0 {
		return nil, errors.New("domains is required")
	}

	access, err := s.repo.GetAccessByIDForUser(ctx, accessID, auth.UserID, auth.Role)
	if err != nil {
		return nil, err
	}

	provider := normalizeOpenProvider(req.Provider)
	accessProvider := normalizeOpenProvider(access.Provider)
	if provider == "" {
		provider = accessProvider
	}
	if provider != accessProvider {
		return nil, fmt.Errorf("provider %s does not match access provider %s", provider, accessProvider)
	}
	if _, ok := providercatalog.OperationDefinition("dns", provider); !ok {
		return nil, fmt.Errorf("unsupported dns provider %s", provider)
	}

	nodeConfig := map[string]interface{}{
		"provider":              provider,
		"accessId":              accessID,
		"domains":               domains,
		"caProvider":            strings.TrimSpace(req.CAProvider),
		"contactEmail":          strings.TrimSpace(req.ContactEmail),
		"keyAlgorithm":          strings.TrimSpace(req.KeyAlgorithm),
		"dnsPropagationTimeout": req.DNSPropagationTimeout,
		"dnsTTL":                req.DNSTTL,
	}

	graph := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{
				"id":   "apply-1",
				"type": "apply",
				"data": map[string]interface{}{
					"name":   "OpenAPI Apply",
					"config": nodeConfig,
				},
			},
			{
				"id":   "end",
				"type": "end",
				"data": map[string]interface{}{
					"name":   "End",
					"config": map[string]interface{}{},
				},
			},
		},
		"edges": []map[string]interface{}{
			{
				"source": "apply-1",
				"target": "end",
			},
		},
	}

	workflow, err := s.repo.SaveWorkflowForUser(ctx, model.Workflow{
		Name:         fmt.Sprintf("OpenAPI Apply %s", time.Now().Format("20060102150405")),
		Description:  "Auto generated by OpenAPI certificate apply API",
		Trigger:      "manual",
		Enabled:      false,
		GraphContent: graph,
		HasContent:   true,
		GraphDraft:   map[string]interface{}{},
		HasDraft:     false,
	}, auth.UserID, auth.Role)
	if err != nil {
		return nil, err
	}

	run, err := s.repo.CreateWorkflowRun(ctx, model.WorkflowRun{
		OwnerUserID: workflow.OwnerUserID,
		WorkflowID:  workflow.ID,
		Status:      "pending",
		Trigger:     "openapi",
		StartedAt:   time.Now(),
		Graph:       workflow.GraphContent,
	})
	if err != nil {
		return nil, err
	}
	s.dispatcher.Start(ctx, run.ID)

	return map[string]interface{}{
		"workflowId": workflow.ID,
		"runId":      run.ID,
		"status":     run.Status,
		"trigger":    run.Trigger,
		"startedAt":  run.StartedAt,
	}, nil
}

func (s *Service) GetOpenCertificateRun(ctx context.Context, auth model.AuthContext, runID string) (map[string]interface{}, error) {
	run, err := s.repo.GetWorkflowRunForUser(ctx, runID, auth.UserID, auth.Role)
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{
		"runId":      run.ID,
		"workflowId": run.WorkflowID,
		"status":     run.Status,
		"trigger":    run.Trigger,
		"startedAt":  run.StartedAt,
		"endedAt":    run.EndedAt,
		"error":      run.Error,
	}

	cert, certErr := s.repo.GetLatestCertificateByRunForUser(ctx, run.ID, auth.UserID, auth.Role)
	if certErr == nil {
		res["certificateId"] = cert.ID
		res["validityNotAfter"] = cert.ValidityNotAfter
	}
	if certErr != nil && !errors.Is(certErr, repository.ErrNotFound) {
		return nil, certErr
	}
	return res, nil
}

func (s *Service) ListOpenCertificateRunEvents(ctx context.Context, auth model.AuthContext, runID, nodeID string, since *time.Time, limit int) ([]model.WorkflowRunEvent, error) {
	return s.repo.ListWorkflowRunEventsForUser(ctx, runID, auth.UserID, auth.Role, nodeID, since, limit)
}
