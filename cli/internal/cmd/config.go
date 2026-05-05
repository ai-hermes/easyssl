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
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current config",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("server: %s\n", cfg.Server)
		if cfg.Token != "" {
			fmt.Println("token: ***")
		} else {
			fmt.Println("token: (none)")
		}
		if cfg.APIKey != "" {
			fmt.Println("api_key: ***")
		} else {
			fmt.Println("api_key: (none)")
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configSetCmd.Long = "Set a config value. Supported keys: server, api_key, token."
	configGetCmd.Long = fmt.Sprintf("Show current config. Default server: %s", config.DefaultServer)
	rootCmd.AddCommand(configCmd)
}
