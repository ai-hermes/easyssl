package cmd

import (
	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/auth/me", nil, nil, client.AuthAuto)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
