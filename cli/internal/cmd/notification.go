package cmd

import (
	"fmt"

	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Notification operations",
}

var notificationTestCmd = &cobra.Command{
	Use:   "test --provider <provider> --access-id <id>",
	Short: "Send a test notification",
	RunE: func(cmd *cobra.Command, args []string) error {
		provider, _ := cmd.Flags().GetString("provider")
		accessID, _ := cmd.Flags().GetString("access-id")
		if provider == "" || accessID == "" {
			return exitErr(2, fmt.Errorf("--provider and --access-id are required"))
		}
		return runRequest("POST", "/api/notifications/test", nil, map[string]any{
			"provider": provider,
			"accessId": accessID,
		}, client.AuthAuto)
	},
}

func init() {
	notificationTestCmd.Flags().String("provider", "", "notification provider")
	notificationTestCmd.Flags().String("access-id", "", "access credential id")
	notificationCmd.AddCommand(notificationTestCmd)
	rootCmd.AddCommand(notificationCmd)
}
