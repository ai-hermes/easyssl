package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows",
}

var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workflows",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Server == "" {
			return fmt.Errorf("not logged in; run 'easyssl login' first")
		}
		c := client.New(cfg)
		data, err := c.ListWorkflows()
		if err != nil {
			return err
		}
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, data, "", "  "); err != nil {
			fmt.Println(string(data))
			return nil
		}
		fmt.Println(pretty.String())
		return nil
	},
}

func init() {
	workflowCmd.AddCommand(workflowListCmd)
	rootCmd.AddCommand(workflowCmd)
}
