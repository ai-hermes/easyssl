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
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object true "email/password"
// @Success 200 {object} util.Response
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	token, admin, err := h.svc.Login(c, req.Email, req.Password)
	if err != nil {
		util.Err(c, http.StatusUnauthorized, err.Error())
		return
	}
	util.OK(c, gin.H{"token": token, "admin": gin.H{"id": admin.ID, "email": admin.Email, "role": admin.Role, "status": admin.Status}})
}

// Me godoc
// @Summary Current user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	auth := h.auth(c)
	admin, err := h.svc.Me(c, auth.UserID)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"id": admin.ID, "email": admin.Email, "role": admin.Role, "status": admin.Status, "authType": auth.AuthType})
}

// ChangePassword godoc
// @Summary Change user password
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param body body object true "new password"
// @Success 200 {object} util.Response
// @Router /auth/password [put]
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

// CreateAPIKey godoc
// @Summary Create API key
// @Tags APIKey
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param body body object true "name/expiresAt"
// @Success 200 {object} util.Response
// @Router /auth/api-keys [post]
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
// @Tags APIKey
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /auth/api-keys [get]
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
// @Tags APIKey
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "API key id"
// @Success 200 {object} util.Response
// @Router /auth/api-keys/{id} [delete]
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
// @Tags OpenAPI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.OpenApplyCertificateRequest true "certificate apply request"
// @Success 200 {object} util.Response
// @Router /open/certificates/apply [post]
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
// @Tags OpenAPI
// @Produce json
// @Security ApiKeyAuth
// @Param runId path string true "Run id"
// @Success 200 {object} util.Response
// @Router /open/certificates/runs/{runId} [get]
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
// @Tags OpenAPI
// @Produce json
// @Security ApiKeyAuth
// @Param runId path string true "Run id"
// @Param nodeId query string false "Node id"
// @Param since query string false "RFC3339 timestamp"
// @Param limit query int false "Limit"
// @Success 200 {object} util.Response
// @Router /open/certificates/runs/{runId}/events [get]
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

// ListProviders godoc
// @Summary List provider definitions
// @Tags Providers
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param kind query string false "provider kind: access, dns, deploy"
// @Success 200 {object} util.Response
// @Router /providers [get]
func (h *Handler) ListProviders(c *gin.Context) {
	items := h.svc.ListProviderDefinitions(c.Query("kind"))
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

// ListAccesses godoc
// @Summary List accesses
// @Tags Access
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /accesses [get]
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
// @Tags Access
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string false "Access id"
// @Param body body object true "access"
// @Success 200 {object} util.Response
// @Router /accesses [post]
// @Router /accesses/{id} [put]
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
// @Tags Access
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Access id"
// @Success 200 {object} util.Response
// @Router /accesses/{id} [delete]
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
// @Tags Access
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Access id"
// @Success 200 {object} util.Response
// @Router /accesses/{id}/test [post]
func (h *Handler) TestAccess(c *gin.Context) {
	if err := h.svc.TestAccess(c, h.auth(c), c.Param("id")); err != nil {
		if err == repository.ErrNotFound {
			util.Err(c, 404, "not found")
			return
		}
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"testedAt": time.Now()})
}

// ListWorkflows godoc
// @Summary List workflows
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /workflows [get]
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
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Success 200 {object} util.Response
// @Router /workflows/{id} [get]
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
// @Tags Workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string false "Workflow id"
// @Param body body object true "workflow"
// @Success 200 {object} util.Response
// @Router /workflows [post]
// @Router /workflows/{id} [put]
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
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Success 200 {object} util.Response
// @Router /workflows/{id} [delete]
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
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Success 200 {object} util.Response
// @Router /workflows/{id}/runs [get]
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
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Param runId path string true "Run id"
// @Success 200 {object} util.Response
// @Router /workflows/{id}/runs/{runId}/nodes [get]
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
// @Router /workflows/{id}/runs/{runId}/events [get]
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
// @Tags Workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Param body body object false "trigger"
// @Success 200 {object} util.Response
// @Router /workflows/{id}/runs [post]
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
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Workflow id"
// @Param runId path string true "Run id"
// @Success 200 {object} util.Response
// @Router /workflows/{id}/runs/{runId}/cancel [post]
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
// @Tags Workflow
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /workflows/stats [get]
func (h *Handler) WorkflowStats(c *gin.Context) { util.OK(c, h.svc.WorkflowStats()) }

// ListCertificates godoc
// @Summary List certificates
// @Tags Certificate
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /certificates [get]
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
// @Tags Certificate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Certificate id"
// @Param body body object false "format"
// @Success 200 {object} util.Response
// @Router /certificates/{id}/download [post]
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
// @Tags Certificate
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param id path string true "Certificate id"
// @Success 200 {object} util.Response
// @Router /certificates/{id}/revoke [post]
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
// @Tags Statistics
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Success 200 {object} util.Response
// @Router /statistics [get]
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
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security ApiKeyAuth
// @Param body body object true "provider/accessId"
// @Success 200 {object} util.Response
// @Router /notifications/test [post]
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
	util.OK(c, gin.H{"sentAt": time.Now()})
}
