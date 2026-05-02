package handler

import (
	"fmt"
	"net/http"
	"time"

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
	util.OK(c, gin.H{"token": token, "admin": gin.H{"id": admin.ID, "email": admin.Email}})
}

func (h *Handler) Me(c *gin.Context) {
	adminID := c.GetString("adminId")
	admin, err := h.svc.Me(c, adminID)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"id": admin.ID, "email": admin.Email})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.ChangePassword(c, c.GetString("adminId"), req.Password); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

func (h *Handler) ListAccesses(c *gin.Context) {
	items, err := h.svc.ListAccesses(c)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

func (h *Handler) SaveAccess(c *gin.Context) {
	var req model.Access
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if id := c.Param("id"); id != "" {
		req.ID = id
	}
	res, err := h.svc.SaveAccess(c, req)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

func (h *Handler) DeleteAccess(c *gin.Context) {
	if err := h.svc.DeleteAccess(c, c.Param("id")); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

func (h *Handler) TestAccess(c *gin.Context) {
	if err := h.svc.TestAccess(c, c.Param("id")); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"testedAt": time.Now()})
}

func (h *Handler) ListWorkflows(c *gin.Context) {
	items, err := h.svc.ListWorkflows(c)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

func (h *Handler) GetWorkflow(c *gin.Context) {
	res, err := h.svc.GetWorkflow(c, c.Param("id"))
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

func (h *Handler) SaveWorkflow(c *gin.Context) {
	var req model.Workflow
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if id := c.Param("id"); id != "" {
		req.ID = id
	}
	res, err := h.svc.SaveWorkflow(c, req)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

func (h *Handler) DeleteWorkflow(c *gin.Context) {
	if err := h.svc.DeleteWorkflow(c, c.Param("id")); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

func (h *Handler) ListWorkflowRuns(c *gin.Context) {
	items, err := h.svc.ListWorkflowRuns(c, c.Param("id"))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

func (h *Handler) ListWorkflowRunNodes(c *gin.Context) {
	items, err := h.svc.ListWorkflowRunNodes(c, c.Param("id"), c.Param("runId"))
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

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
	items, err := h.svc.ListWorkflowRunEvents(c, c.Param("id"), c.Param("runId"), c.Query("nodeId"), since, limit)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

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
	runID, err := h.svc.StartWorkflowRun(c, c.Param("id"), req.Trigger)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"runId": runID})
}

func (h *Handler) CancelWorkflowRun(c *gin.Context) {
	if err := h.svc.CancelWorkflowRun(c, c.Param("runId")); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

func (h *Handler) WorkflowStats(c *gin.Context) { util.OK(c, h.svc.WorkflowStats()) }

func (h *Handler) ListCertificates(c *gin.Context) {
	items, err := h.svc.ListCertificates(c)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"items": items, "totalItems": len(items)})
}

func (h *Handler) DownloadCertificate(c *gin.Context) {
	var req struct {
		Format string `json:"format"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.Format == "" {
		req.Format = "PEM"
	}
	res, err := h.svc.DownloadCertificate(c, c.Param("id"), req.Format)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

func (h *Handler) RevokeCertificate(c *gin.Context) {
	if err := h.svc.RevokeCertificate(c, c.Param("id")); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{})
}

func (h *Handler) Statistics(c *gin.Context) {
	res, err := h.svc.Statistics(c)
	if err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, res)
}

func (h *Handler) TestNotification(c *gin.Context) {
	var req struct {
		Provider string `json:"provider"`
		AccessID string `json:"accessId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.TestNotification(c, req.Provider, req.AccessID); err != nil {
		util.Err(c, 500, err.Error())
		return
	}
	util.OK(c, gin.H{"sentAt": time.Now()})
}
