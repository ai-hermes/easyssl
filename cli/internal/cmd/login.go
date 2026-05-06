package cmd

import (
	"fmt"

	"easyssl/cli/internal/client"
	"easyssl/cli/internal/config"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the EasySSL server",
	Long: `Login to an EasySSL server using an API key.
The API key is persisted to local config for subsequent commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server, _ := cmd.Flags().GetString("server")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if server == "" {
			server = cfg.Server
		}
		if apiKey == "" {
			return exitErr(2, fmt.Errorf("--api-key is required"))
		}

		next := cfg
		next.Server = server
		next.APIKey = apiKey
		next.Token = ""

		c, err := client.New(next, client.Options{})
		if err != nil {
			return exitErr(2, err)
		}
		if _, err := c.Do("GET", "/api/auth/me", nil, nil, client.AuthAuto); err != nil {
			return parseAPIError(fmt.Errorf("login failed: %w", err))
		}
		if err := config.Save(next); err != nil {
			return exitErr(5, fmt.Errorf("save config: %w", err))
		}
		cfg = next
		return printOutput(map[string]any{"ok": true, "msg": "login successful"})
	},
}

func init() {
	loginCmd.Flags().String("server", "", "EasySSL server URL")
	loginCmd.Flags().String("api-key", "", "EasySSL API key")
	rootCmd.AddCommand(loginCmd)
}
