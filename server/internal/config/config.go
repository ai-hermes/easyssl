package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr  string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	// Load order: root .env -> server .env -> process env.
	_ = godotenv.Load("../.env", ".env")

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		pgHost := getEnv("PG_HOST", "127.0.0.1")
		pgPort := getEnv("PG_PORT", "5432")
		pgUser := getEnv("PG_USER", "postgres")
		pgPassword := getEnv("PG_PASSWORD", "postgres")
		pgDatabase := getEnv("PG_DATABASE", "easyssl")
		pgSSLMode := getEnv("PG_SSLMODE", "disable")
		databaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			pgUser,
			pgPassword,
			pgHost,
			pgPort,
			pgDatabase,
			pgSSLMode,
		)
	}

	return Config{
		ListenAddr:  getEnv("LISTEN_ADDR", ":8090"),
		DatabaseURL: databaseURL,
		JWTSecret:   getEnv("JWT_SECRET", "easyssl-dev-secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
