package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr  string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		ListenAddr:  getEnv("LISTEN_ADDR", ":8090"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/easyssl?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "easyssl-dev-secret"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
