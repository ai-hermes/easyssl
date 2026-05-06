package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const DefaultServer = "https://easyssl.spotty.com.cn"

// Config holds CLI configuration and credentials.
type Config struct {
	Server  string `mapstructure:"server"`
	Token   string `mapstructure:"token"`
	APIKey  string `mapstructure:"api_key"`
	Output  string `mapstructure:"output"`
	Timeout int    `mapstructure:"timeout"`
	Trace   bool   `mapstructure:"trace"`
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
	cfg := Config{
		Server:  DefaultServer,
		Output:  "json",
		Timeout: 30,
	}
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
	if cfg.Server == "" {
		cfg.Server = DefaultServer
	}
	if cfg.Output == "" {
		cfg.Output = "json"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30
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
	v.Set("output", cfg.Output)
	v.Set("timeout", cfg.Timeout)
	v.Set("trace", cfg.Trace)
	if _, err := os.Stat(File()); os.IsNotExist(err) {
		return v.WriteConfigAs(File())
	}
	return v.WriteConfig()
}

// Clear removes the stored config file.
func Clear() error {
	return os.RemoveAll(File())
}
