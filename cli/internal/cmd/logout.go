package cmd

import (
	"fmt"

	"easyssl/cli/internal/config"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Clear(); err != nil {
			return exitErr(5, fmt.Errorf("clear config: %w", err))
		}
		return printOutput(map[string]any{"ok": true, "msg": "logged out"})
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
