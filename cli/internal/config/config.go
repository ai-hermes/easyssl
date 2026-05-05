package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds CLI configuration and credentials.
type Config struct {
	Server   string `mapstructure:"server"`
	Token    string `mapstructure:"token"`
	APIKey   string `mapstructure:"api_key"`
}

// Dir returns the CLI config directory.
func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "easyssl")
}

// File returns the full path to the config file.
func File() string {
	return filepath.Join(Dir(), "config.yaml")
}

// Load reads the config file into a Config value.
func Load() (Config, error) {
	var cfg Config
	v := viper.New()
	v.SetConfigFile(File())
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config: %w", err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}
	return cfg, nil
}

// Save persists the config to disk.
func Save(cfg Config) error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	v := viper.New()
	v.SetConfigFile(File())
	v.SetConfigType("yaml")
	_ = v.ReadInConfig()
	v.Set("server", cfg.Server)
	v.Set("token", cfg.Token)
	v.Set("api_key", cfg.APIKey)
	return v.WriteConfig()
}

// Clear removes the stored config file.
func Clear() error {
	return os.RemoveAll(File())
}
