package handler

import (
	"fmt"
	"net/http"
	"time"

	"easyssl/server/internal/middleware"
	"easyssl/server/internal/model"
	"easyssl/server/internal/repository"
	"easyssl/server/internal/service"
	"easyssl/server/internal/util"
	"easyssl/server/internal/version"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) auth(c *gin.Context) model.AuthContext {
	return middleware.GetAuthContext(c)
}

// Login godoc
// @Summary User login
// @Description Authenticate with email and password to obtain a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "email/password"
// @Success 200 {object} util.Response
// @Router /api/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	token, user, err := h.svc.Login(c, req.Email, req.Password)
	if err != nil {
		util.Err(c, http.StatusUnauthorized, err.Error())
		return
	}
	util.OK(c, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email, "role": user.Role, "status": user.Status}})
}

// Me godoc
// @Summary Current user
// @Description Get the current authenticated user's profile information.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	auth := h.auth(c)
	user, err := h.svc.Me(c, auth.UserID)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"id": user.ID, "email": user.Email, "role": user.Role, "status": user.Status, "authType": auth.AuthType})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Update the password for the current user. Takes effect on next login.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param body body model.ChangePasswordRequest true "new password"
// @Success 200 {object} util.Response
// @Router /api/auth/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.ChangePassword(c, h.auth(c).UserID, req.Password); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

// Register godoc
// @Summary User registration
// @Description Register a new user account with email and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.RegisterRequest true "email/password"
// @Success 200 {object} util.Response
// @Router /api/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	user, err := h.svc.Register(c, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already registered" {
			util.Err(c, 409, err.Error())
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"id": user.ID, "email": user.Email, "role": user.Role, "status": user.Status})
}

// CreateAPIKey godoc
// @Summary Create API key
// @Description Create a new API key for programmatic access via X-API-Key header.
// @Tags APIKey
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param body body model.CreateAPIKeyRequest true "name/expiresAt"
// @Success 200 {object} util.Response
// @Router /api/auth/api-keys [post]
func (h *Handler) CreateAPIKey(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ExpiresAt string `json:"expiresAt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			util.Err(c, 400, "invalid expiresAt, must be RFC3339")
			return
		}
		expiresAt = &t
	}
	res, err := h.svc.CreateAPIKey(c, h.auth(c), req.Name, expiresAt)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// ListAPIKeys godoc
// @Summary List API keys
// @Description List all API keys created by the current user.
// @Tags APIKey
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/auth/api-keys [get]
func (h *Handler) ListAPIKeys(c *gin.Context) {
	items, err := h.svc.ListAPIKeys(c, h.auth(c))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// RevokeAPIKey godoc
// @Summary Revoke API key
// @Description Revoke an existing API key to disable its access.
// @Tags APIKey
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "API key id"
// @Success 200 {object} util.Response
// @Router /api/auth/api-keys/{id} [delete]
func (h *Handler) RevokeAPIKey(c *gin.Context) {
	if err := h.svc.RevokeAPIKey(c, h.auth(c), c.Param("id")); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

// OpenApplyCertificate godoc
// @Summary Apply certificate via OpenAPI
// @Description Apply for an SSL certificate using the OpenAPI endpoint with API key authentication.
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.OpenApplyCertificateRequest true "certificate apply request"
// @Success 200 {object} util.Response
// @Router /openapi/certificates/apply [post]
func (h *Handler) OpenApplyCertificate(c *gin.Context) {
	var req model.OpenApplyCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	res, err := h.svc.OpenApplyCertificate(c, h.auth(c), req)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 400, err.Error())
		return
	}
	util.OK(c, res)
}

// GetOpenCertificateRun godoc
// @Summary Get OpenAPI certificate apply run
// @Description Get the status and result of a certificate apply run triggered via OpenAPI.
// @Tags OpenAPI
// @Produce json
// @Security ApiKeyAuth
// @Param runId path string true "Run id"
// @Success 200 {object} util.Response
// @Router /openapi/certificates/runs/{runId} [get]
func (h *Handler) GetOpenCertificateRun(c *gin.Context) {
	res, err := h.svc.GetOpenCertificateRun(c, h.auth(c), c.Param("runId"))
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// ListOpenCertificateRunEvents godoc
// @Summary List OpenAPI certificate apply run events
// @Description List execution events for a certificate apply run.
// @Tags OpenAPI
// @Produce json
// @Security ApiKeyAuth
// @Param runId path string true "Run id"
// @Param nodeId query string false "Node id"
// @Param since query string false "RFC3339 timestamp"
// @Param limit query int false "Limit"
// @Success 200 {object} util.Response
// @Router /openapi/certificates/runs/{runId}/events [get]
func (h *Handler) ListOpenCertificateRunEvents(c *gin.Context) {
	var since *time.Time
	if raw := c.Query("since"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			util.Err(c, 400, "invalid since, must be RFC3339")
			return
		}
		since = &t
	}

	limit := 200
	if rawLimit := c.Query("limit"); rawLimit != "" {
		var parsed int
		if _, err := fmt.Sscanf(rawLimit, "%d", &parsed); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	items, err := h.svc.ListOpenCertificateRunEvents(c, h.auth(c), c.Param("runId"), c.Query("nodeId"), since, limit)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// OpenListAccesses godoc
// @Summary List accesses (OpenAPI)
// @Description List all access credentials configured by the current user via OpenAPI.
// @Tags OpenAPI
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /openapi/accesses [get]
func (h *Handler) OpenListAccesses(c *gin.Context) {
	items, err := h.svc.ListAccesses(c, h.auth(c))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// OpenListCertificates godoc
// @Summary List certificates (OpenAPI)
// @Description List all certificates managed by the current user via OpenAPI.
// @Tags OpenAPI
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /openapi/certificates [get]
func (h *Handler) OpenListCertificates(c *gin.Context) {
	items, err := h.svc.ListCertificates(c, h.auth(c))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// OpenDownloadCertificate godoc
// @Summary Download certificate (OpenAPI)
// @Description Download a certificate in the specified format (PEM, PFX, JKS) via OpenAPI.
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Certificate id"
// @Param body body model.DownloadCertificateRequest false "format"
// @Success 200 {object} util.Response
// @Router /openapi/certificates/{id}/download [post]
func (h *Handler) OpenDownloadCertificate(c *gin.Context) {
	var req struct {
		Format string `json:"format"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.Format == "" {
		req.Format = "PEM"
	}
	res, err := h.svc.DownloadCertificate(c, h.auth(c), c.Param("id"), req.Format)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// ListProviders godoc
// @Summary List provider definitions
// @Description List all available provider definitions for DNS, access, and deploy operations.
// @Tags Providers
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param kind query string false "provider kind: access, dns, deploy"
// @Success 200 {object} util.Response
// @Router /api/providers [get]
func (h *Handler) ListProviders(c *gin.Context) {
	items := h.svc.ListProviderDefinitions(c.Query("kind"))
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// ListAccesses godoc
// @Summary List accesses
// @Description List all access credentials configured by the current user.
// @Tags Access
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/accesses [get]
func (h *Handler) ListAccesses(c *gin.Context) {
	items, err := h.svc.ListAccesses(c, h.auth(c))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// SaveAccess godoc
// @Summary Create or update access
// @Description Create a new access credential or update an existing one.
// @Tags Access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string false "Access id"
// @Param body body model.Access true "access"
// @Success 200 {object} util.Response
// @Router /api/accesses [post]
// @Router /api/accesses/{id} [put]
func (h *Handler) SaveAccess(c *gin.Context) {
	var req model.Access
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if id := c.Param("id"); id != "" {
		req.ID = id
	}
	res, err := h.svc.SaveAccess(c, h.auth(c), req)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// DeleteAccess godoc
// @Summary Delete access
// @Description Soft delete an access credential.
// @Tags Access
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Access id"
// @Success 200 {object} util.Response
// @Router /api/accesses/{id} [delete]
func (h *Handler) DeleteAccess(c *gin.Context) {
	if err := h.svc.DeleteAccess(c, h.auth(c), c.Param("id")); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

// TestAccess godoc
// @Summary Test access
// @Description Test the connectivity and validity of an access credential.
// @Tags Access
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Access id"
// @Success 200 {object} util.Response
// @Router /api/accesses/{id}/test [post]
func (h *Handler) TestAccess(c *gin.Context) {
	if err := h.svc.TestAccess(c, h.auth(c), c.Param("id")); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"testedAt": time.Now().Format("2006-01-02 15:04:05")})
}

// ListWorkflows godoc
// @Summary List workflows
// @Description List all workflows owned by the current user.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/workflows [get]
func (h *Handler) ListWorkflows(c *gin.Context) {
	items, err := h.svc.ListWorkflows(c, h.auth(c))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// GetWorkflow godoc
// @Summary Get workflow
// @Description Get detailed information of a specific workflow.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Success 200 {object} util.Response
// @Router /api/workflows/{id} [get]
func (h *Handler) GetWorkflow(c *gin.Context) {
	res, err := h.svc.GetWorkflow(c, h.auth(c), c.Param("id"))
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// SaveWorkflow godoc
// @Summary Create or update workflow
// @Description Create a new workflow or update an existing workflow definition.
// @Tags Workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string false "Workflow id"
// @Param body body model.Workflow true "workflow"
// @Success 200 {object} util.Response
// @Router /api/workflows [post]
// @Router /api/workflows/{id} [put]
func (h *Handler) SaveWorkflow(c *gin.Context) {
	var req model.Workflow
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if id := c.Param("id"); id != "" {
		req.ID = id
	}
	res, err := h.svc.SaveWorkflow(c, h.auth(c), req)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// DeleteWorkflow godoc
// @Summary Delete workflow
// @Description Delete a workflow and all its associated data.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Success 200 {object} util.Response
// @Router /api/workflows/{id} [delete]
func (h *Handler) DeleteWorkflow(c *gin.Context) {
	if err := h.svc.DeleteWorkflow(c, h.auth(c), c.Param("id")); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

// ListWorkflowRuns godoc
// @Summary List workflow runs
// @Description List execution runs for a specific workflow.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Success 200 {object} util.Response
// @Router /api/workflows/{id}/runs [get]
func (h *Handler) ListWorkflowRuns(c *gin.Context) {
	items, err := h.svc.ListWorkflowRuns(c, h.auth(c), c.Param("id"))
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// ListWorkflowRunNodes godoc
// @Summary List workflow run nodes
// @Description List all node execution records for a workflow run.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Param runId path string true "Run id"
// @Success 200 {object} util.Response
// @Router /api/workflows/{id}/runs/{runId}/nodes [get]
func (h *Handler) ListWorkflowRunNodes(c *gin.Context) {
	items, err := h.svc.ListWorkflowRunNodes(c, h.auth(c), c.Param("id"), c.Param("runId"))
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// ListWorkflowRunEvents godoc
// @Summary List workflow run events
// @Description List event logs for a workflow run, optionally filtered by node.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Param runId path string true "Run id"
// @Param nodeId query string false "Node id"
// @Param since query string false "RFC3339 timestamp"
// @Param limit query int false "Limit"
// @Success 200 {object} util.Response
// @Router /api/workflows/{id}/runs/{runId}/events [get]
func (h *Handler) ListWorkflowRunEvents(c *gin.Context) {
	var since *time.Time
	if raw := c.Query("since"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			util.Err(c, 400, "invalid since, must be RFC3339")
			return
		}
		since = &t
	}
	limit := 200
	if rawLimit := c.Query("limit"); rawLimit != "" {
		var parsed int
		if _, err := fmt.Sscanf(rawLimit, "%d", &parsed); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	items, err := h.svc.ListWorkflowRunEvents(c, h.auth(c), c.Param("id"), c.Param("runId"), c.Query("nodeId"), since, limit)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// StartWorkflowRun godoc
// @Summary Start workflow run
// @Description Manually trigger a workflow execution.
// @Tags Workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Param body body model.StartWorkflowRunRequest false "trigger"
// @Success 200 {object} util.Response
// @Router /api/workflows/{id}/runs [post]
func (h *Handler) StartWorkflowRun(c *gin.Context) {
	var req struct {
		Trigger string `json:"trigger"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if req.Trigger == "" {
		req.Trigger = "manual"
	}
	runID, err := h.svc.StartWorkflowRun(c, h.auth(c), c.Param("id"), req.Trigger)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"runId": runID})
}

// CancelWorkflowRun godoc
// @Summary Cancel workflow run
// @Description Cancel a running or pending workflow execution.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Param runId path string true "Run id"
// @Success 200 {object} util.Response
// @Router /api/workflows/{id}/runs/{runId}/cancel [post]
func (h *Handler) CancelWorkflowRun(c *gin.Context) {
	if err := h.svc.CancelWorkflowRun(c, h.auth(c), c.Param("runId")); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

// WorkflowStats godoc
// @Summary Workflow dispatcher stats
// @Description Get real-time workflow dispatcher statistics including concurrency and queue status.
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/workflows/stats [get]
func (h *Handler) WorkflowStats(c *gin.Context) { util.OK(c, h.svc.WorkflowStats()) }

// ListCertificates godoc
// @Summary List certificates
// @Description List all certificates managed by the current user.
// @Tags Certificate
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/certificates [get]
func (h *Handler) ListCertificates(c *gin.Context) {
	items, err := h.svc.ListCertificates(c, h.auth(c))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// DownloadCertificate godoc
// @Summary Download certificate
// @Description Download a certificate in the specified format (PEM, PFX, JKS).
// @Tags Certificate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Certificate id"
// @Param body body model.DownloadCertificateRequest false "format"
// @Success 200 {object} util.Response
// @Router /api/certificates/{id}/download [post]
func (h *Handler) DownloadCertificate(c *gin.Context) {
	var req struct {
		Format string `json:"format"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.Format == "" {
		req.Format = "PEM"
	}
	res, err := h.svc.DownloadCertificate(c, h.auth(c), c.Param("id"), req.Format)
	if err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// RevokeCertificate godoc
// @Summary Revoke certificate
// @Description Revoke a certificate to mark it as invalid.
// @Tags Certificate
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Certificate id"
// @Success 200 {object} util.Response
// @Router /api/certificates/{id}/revoke [post]
func (h *Handler) RevokeCertificate(c *gin.Context) {
	if err := h.svc.RevokeCertificate(c, h.auth(c), c.Param("id")); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

// Statistics godoc
// @Summary Dashboard statistics
// @Description Get dashboard statistics including certificate and workflow counts.
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/statistics [get]
func (h *Handler) Statistics(c *gin.Context) {
	res, err := h.svc.Statistics(c, h.auth(c))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

// TestNotification godoc
// @Summary Test notification
// @Description Send a test notification using the specified provider and access.
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param body body model.TestNotificationRequest true "provider/accessId"
// @Success 200 {object} util.Response
// @Router /api/notifications/test [post]
func (h *Handler) TestNotification(c *gin.Context) {
	var req struct {
		Provider string `json:"provider"`
		AccessID string `json:"accessId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.TestNotification(c, h.auth(c), req.Provider, req.AccessID); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"sentAt": time.Now().Format("2006-01-02 15:04:05")})
}

// ListUsers godoc
// @Summary List all users (admin only)
// @Description List all registered users. Requires admin role.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/admin/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	if !h.auth(c).IsAdmin() {
		util.Err(c, 403, "admin role required")
		return
	}
	items, err := h.svc.ListUsers(c)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	out := make([]gin.H, 0, len(items))
	for _, u := range items {
		out = append(out, gin.H{"id": u.ID, "email": u.Email, "role": u.Role, "status": u.Status, "createdAt": u.CreatedAt})
	}
	util.OK(c, gin.H{"items": out, "totalItems": len(out)})
}

// UpdateUserStatus godoc
// @Summary Update user status (admin only)
// @Description Enable or disable a user account. Requires admin role.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "User id"
// @Param body body map[string]string true "{status}"
// @Success 200 {object} util.Response
// @Router /api/admin/users/{id}/status [put]
func (h *Handler) UpdateUserStatus(c *gin.Context) {
	if !h.auth(c).IsAdmin() {
		util.Err(c, 403, "admin role required")
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateUserStatus(c, c.Param("id"), req.Status); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

// Version godoc
// @Summary Get server version (admin only)
// @Description Get the current server build version including git branch and commit id.
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /api/admin/version [get]
func (h *Handler) Version(c *gin.Context) {
	if !h.auth(c).IsAdmin() {
		util.Err(c, 403, "admin role required")
		return
	}
	util.OK(c, gin.H{"version": version.String(), "commitUrl": version.CommitURL()})
}
