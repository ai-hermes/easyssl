package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSave(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("HOME", oldHome)
	}()
	_ = os.Setenv("HOME", tmpDir)

	cfg := Config{
		Server:  "http://localhost:8090",
		Token:   "test-token",
		APIKey:  "test-api-key",
		Output:  "json",
		Timeout: 42,
		Trace:   true,
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if loaded.Server != cfg.Server {
		t.Errorf("server = %q, want %q", loaded.Server, cfg.Server)
	}
	if loaded.Token != cfg.Token {
		t.Errorf("token = %q, want %q", loaded.Token, cfg.Token)
	}
	if loaded.APIKey != cfg.APIKey {
		t.Errorf("api_key = %q, want %q", loaded.APIKey, cfg.APIKey)
	}
	if loaded.Output != cfg.Output {
		t.Errorf("output = %q, want %q", loaded.Output, cfg.Output)
	}
	if loaded.Timeout != cfg.Timeout {
		t.Errorf("timeout = %d, want %d", loaded.Timeout, cfg.Timeout)
	}
	if loaded.Trace != cfg.Trace {
		t.Errorf("trace = %v, want %v", loaded.Trace, cfg.Trace)
	}
}

func TestClear(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("HOME", oldHome)
	}()
	_ = os.Setenv("HOME", tmpDir)

	cfg := Config{Server: "http://localhost:8090"}
	if err := Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	if err := Clear(); err != nil {
		t.Fatalf("clear config: %v", err)
	}
	if _, err := os.Stat(filepath.Join(Dir(), "config.yaml")); !os.IsNotExist(err) {
		t.Error("config file should not exist after clear")
	}
}

func TestLoadUsesDefaultsWhenConfigMissing(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() {
		_ = os.Setenv("HOME", oldHome)
	}()
	_ = os.Setenv("HOME", tmpDir)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if loaded.Server != DefaultServer {
		t.Fatalf("server = %q, want %q", loaded.Server, DefaultServer)
	}
	if loaded.Output != "json" {
		t.Fatalf("output = %q, want %q", loaded.Output, "json")
	}
	if loaded.Timeout != 30 {
		t.Fatalf("timeout = %d, want 30", loaded.Timeout)
	}
}
