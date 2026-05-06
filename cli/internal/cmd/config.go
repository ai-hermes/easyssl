package cmd

import (
	"fmt"

	"easyssl/cli/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]
		switch key {
		case "server":
			cfg.Server = value
		case "api_key":
			cfg.APIKey = value
		case "token":
			cfg.Token = value
		case "output":
			if !isOutputFormatSupported(value) {
				return exitErr(2, fmt.Errorf("unsupported output: %s", value))
			}
			cfg.Output = value
		case "timeout":
			var parsed int
			if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil || parsed <= 0 {
				return exitErr(2, fmt.Errorf("timeout must be positive integer seconds"))
			}
			cfg.Timeout = parsed
		case "trace":
			cfg.Trace = value == "true" || value == "1"
		default:
			return exitErr(2, fmt.Errorf("unknown config key: %s", key))
		}
		if err := config.Save(cfg); err != nil {
			return exitErr(5, fmt.Errorf("save config: %w", err))
		}
		return printOutput(map[string]any{"ok": true, "key": key, "value": value})
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printOutput(map[string]any{
			"server":  cfg.Server,
			"token":   mask(cfg.Token),
			"api_key": mask(cfg.APIKey),
			"output":  cfg.Output,
			"timeout": cfg.Timeout,
			"trace":   cfg.Trace,
		})
	},
}

func mask(s string) string {
	if s == "" {
		return ""
	}
	return "***"
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configSetCmd.Long = "Set a config value. Supported keys: server, api_key, token, output, timeout, trace."
	configGetCmd.Long = fmt.Sprintf("Show current config. Default server: %s", config.DefaultServer)
	rootCmd.AddCommand(configCmd)
}
