package cmd

import (
	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var statisticsCmd = &cobra.Command{
	Use:   "statistics",
	Short: "Dashboard statistics",
}

var statisticsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/statistics", nil, nil, client.AuthAuto)
	},
}

func init() {
	statisticsCmd.AddCommand(statisticsGetCmd)
	rootCmd.AddCommand(statisticsCmd)
}
