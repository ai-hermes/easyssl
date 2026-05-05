package router

import (
	"context"
	"net/http"
	"os"

	"easyssl/server/internal/config"
	"easyssl/server/internal/db"
	"easyssl/server/internal/handler"
	"easyssl/server/internal/middleware"
	"easyssl/server/internal/repository"
	"easyssl/server/internal/service"
	"easyssl/server/internal/workflow"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config, database *db.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	repo := repository.New(database)
	if err := repo.EnsureWorkflowRunTables(context.Background()); err != nil {
		panic(err)
	}
	dispatcher := workflow.NewDispatcher(repo, 2)
	svc := service.New(repo, cfg.JWTSecret, dispatcher)
	_ = svc.EnsureBootstrapUser(context.Background(), "admin@easyssl.local", "1234567890")
	h := handler.New(svc)

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.GET("/api/openapi.json", func(c *gin.Context) {
		c.File("docs/swagger.json")
	})

	auth := r.Group("/api/auth")
	auth.POST("/login", h.Login)
	auth.POST("/register", h.Register)

	api := r.Group("/api")
	api.Use(middleware.RequireAuth(cfg.JWTSecret, svc))
	api.GET("/auth/me", h.Me)
	api.PUT("/auth/password", h.ChangePassword)
	api.POST("/auth/api-keys", h.CreateAPIKey)
	api.GET("/auth/api-keys", h.ListAPIKeys)
	api.DELETE("/auth/api-keys/:id", h.RevokeAPIKey)

	api.GET("/admin/users", h.ListUsers)
	api.PUT("/admin/users/:id/status", h.UpdateUserStatus)
	api.GET("/admin/version", h.Version)

	api.GET("/providers", h.ListProviders)

	api.GET("/accesses", h.ListAccesses)
	api.POST("/accesses", h.SaveAccess)
	api.PUT("/accesses/:id", h.SaveAccess)
	api.DELETE("/accesses/:id", h.DeleteAccess)
	api.POST("/accesses/:id/test", h.TestAccess)

	api.GET("/workflows", h.ListWorkflows)
	api.POST("/workflows", h.SaveWorkflow)
	api.GET("/workflows/:id", h.GetWorkflow)
	api.PUT("/workflows/:id", h.SaveWorkflow)
	api.DELETE("/workflows/:id", h.DeleteWorkflow)
	api.GET("/workflows/:id/runs", h.ListWorkflowRuns)
	api.GET("/workflows/:id/runs/:runId/nodes", h.ListWorkflowRunNodes)
	api.GET("/workflows/:id/runs/:runId/events", h.ListWorkflowRunEvents)
	api.POST("/workflows/:id/runs", h.StartWorkflowRun)
	api.POST("/workflows/:id/runs/:runId/cancel", h.CancelWorkflowRun)
	api.GET("/workflows/stats", h.WorkflowStats)

	api.GET("/certificates", h.ListCertificates)
	api.POST("/certificates/:id/download", h.DownloadCertificate)
	api.POST("/certificates/:id/revoke", h.RevokeCertificate)

	api.GET("/statistics", h.Statistics)
	api.POST("/notifications/test", h.TestNotification)

	open := r.Group("/openapi")
	open.Use(middleware.RequireAPIKeyAuth(svc))
	open.GET("/accesses", h.OpenListAccesses)
	open.GET("/certificates", h.OpenListCertificates)
	open.POST("/certificates/:id/download", h.OpenDownloadCertificate)
	open.POST("/certificates/apply", h.OpenApplyCertificate)
	open.GET("/certificates/runs/:runId", h.GetOpenCertificateRun)
	open.GET("/certificates/runs/:runId/events", h.ListOpenCertificateRunEvents)

	if _, err := os.Stat("../web/dist/index.html"); err == nil {
		r.Static("/assets", "../web/dist/assets")
		r.StaticFile("/favicon.png", "../web/dist/favicon.png")
		r.StaticFile("/logo.png", "../web/dist/logo.png")
		r.StaticFile("/logo-v1.png", "../web/dist/logo-v1.png")
		r.StaticFile("/logo-v2.png", "../web/dist/logo-v2.png")
		r.NoRoute(func(c *gin.Context) { c.File("../web/dist/index.html") })
	}

	return r
}
