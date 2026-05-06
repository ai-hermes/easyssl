package cmd

import (
	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Provider definitions",
}

var providerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		kind, _ := cmd.Flags().GetString("kind")
		return runRequest("GET", "/api/providers", map[string]string{"kind": kind}, nil, client.AuthAuto)
	},
}

func init() {
	providerListCmd.Flags().String("kind", "", "provider kind: access|dns|deploy")
	providerCmd.AddCommand(providerListCmd)
	rootCmd.AddCommand(providerCmd)
}
