package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"easyssl/server/internal/config"
	"easyssl/server/internal/db"
)

func main() {
	cfg := config.Load()
	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	ctx := context.Background()
	if _, err := database.Pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version BIGINT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		panic(err)
	}

	entries, err := os.ReadDir("migrations")
	if err != nil {
		panic(err)
	}

	files := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, file := range files {
		parts := strings.SplitN(file, "_", 2)
		v, _ := strconv.ParseInt(parts[0], 10, 64)
		var exists int
		err := database.Pool.QueryRow(ctx, `SELECT COUNT(1) FROM schema_migrations WHERE version=$1`, v).Scan(&exists)
		if err != nil {
			panic(err)
		}
		if exists > 0 {
			continue
		}

		b, err := os.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			panic(err)
		}

		tx, err := database.Pool.Begin(ctx)
		if err != nil {
			panic(err)
		}
		if _, err := tx.Exec(ctx, string(b)); err != nil {
			_ = tx.Rollback(ctx)
			panic(fmt.Errorf("migration %s failed: %w", file, err))
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES($1)`, v); err != nil {
			_ = tx.Rollback(ctx)
			panic(err)
		}
		if err := tx.Commit(ctx); err != nil {
			panic(err)
		}
		fmt.Printf("applied migration: %s\n", file)
	}
}
