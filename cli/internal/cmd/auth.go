package cmd

import (
	"fmt"

	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication and API key operations",
}

var authWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/auth/me", nil, nil, client.AuthAuto)
	},
}

var authPasswordCmd = &cobra.Command{
	Use:   "password <new-password>",
	Short: "Change account password",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("PUT", "/api/auth/password", nil, map[string]any{"password": args[0]}, client.AuthAuto)
	},
}

var apiKeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Manage API keys",
}

var apiKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/auth/api-keys", nil, nil, client.AuthAuto)
	},
}

var apiKeyCreateCmd = &cobra.Command{
	Use:   "create --name <name>",
	Short: "Create an API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		expiresAt, _ := cmd.Flags().GetString("expires-at")
		if name == "" {
			return exitErr(2, fmt.Errorf("--name is required"))
		}
		body := map[string]any{"name": name}
		if expiresAt != "" {
			body["expiresAt"] = expiresAt
		}
		return runRequest("POST", "/api/auth/api-keys", nil, body, client.AuthAuto)
	},
}

var apiKeyRevokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Revoke an API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("DELETE", "/api/auth/api-keys/"+args[0], nil, nil, client.AuthAuto)
	},
}

func init() {
	apiKeyCreateCmd.Flags().String("name", "", "API key name")
	apiKeyCreateCmd.Flags().String("expires-at", "", "RFC3339 expiration timestamp")

	authCmd.AddCommand(authWhoamiCmd)
	authCmd.AddCommand(authPasswordCmd)
	apiKeyCmd.AddCommand(apiKeyListCmd)
	apiKeyCmd.AddCommand(apiKeyCreateCmd)
	apiKeyCmd.AddCommand(apiKeyRevokeCmd)
	authCmd.AddCommand(apiKeyCmd)
	rootCmd.AddCommand(authCmd)
}
