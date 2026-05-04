package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"easyssl/server/internal/config"
	"easyssl/server/internal/db"
	"easyssl/server/internal/router"
)

// @title EasySSL API
// @version 1.0
// @description EasySSL service APIs.
// @description 1) Call POST /api/auth/login to get JWT token.
// @description 2) Call POST /api/auth/api-keys to generate API key token (shown once).
// @description 3) Call OpenAPI certificate apply endpoints with `X-API-Key: <token>`.
// @description 4) Other business APIs support `Authorization: Bearer <token>` or `X-API-Key: <token>`.
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: Bearer <token>
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	cfg := config.Load()

	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer database.Close()

	r := router.New(cfg, database)
	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("easyssl api listening at %s", cfg.ListenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped: %v", err)
	}
}
