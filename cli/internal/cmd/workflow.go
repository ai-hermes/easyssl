package cmd

import (
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
		return runRequest("GET", "/api/workflows", nil, nil, client.AuthAuto)
	},
}

var workflowGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get workflow details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/workflows/"+args[0], nil, nil, client.AuthAuto)
	},
}

var workflowCreateCmd = &cobra.Command{
	Use:   "create --file <json> | --data <json>",
	Short: "Create workflow",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		data, _ := cmd.Flags().GetString("data")
		body, err := readJSONInput(file, data)
		if err != nil {
			return exitErr(2, err)
		}
		return runRequest("POST", "/api/workflows", nil, body, client.AuthAuto)
	},
}

var workflowUpdateCmd = &cobra.Command{
	Use:   "update <id> --file <json> | --data <json>",
	Short: "Update workflow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		data, _ := cmd.Flags().GetString("data")
		body, err := readJSONInput(file, data)
		if err != nil {
			return exitErr(2, err)
		}
		return runRequest("PUT", "/api/workflows/"+args[0], nil, body, client.AuthAuto)
	},
}

var workflowDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete workflow",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("DELETE", "/api/workflows/"+args[0], nil, nil, client.AuthAuto)
	},
}

var workflowStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Workflow dispatcher stats",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/workflows/stats", nil, nil, client.AuthAuto)
	},
}

var workflowRunCmd = &cobra.Command{Use: "run", Short: "Manage workflow runs"}

var workflowRunListCmd = &cobra.Command{
	Use:   "list <workflow-id>",
	Short: "List workflow runs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/workflows/"+args[0]+"/runs", nil, nil, client.AuthAuto)
	},
}

var workflowRunStartCmd = &cobra.Command{
	Use:   "start <workflow-id>",
	Short: "Start a workflow run",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		trigger, _ := cmd.Flags().GetString("trigger")
		body := map[string]any{}
		if trigger != "" {
			body["trigger"] = trigger
		}
		return runRequest("POST", "/api/workflows/"+args[0]+"/runs", nil, body, client.AuthAuto)
	},
}

var workflowRunCancelCmd = &cobra.Command{
	Use:   "cancel <workflow-id> <run-id>",
	Short: "Cancel workflow run",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("POST", "/api/workflows/"+args[0]+"/runs/"+args[1]+"/cancel", nil, nil, client.AuthAuto)
	},
}

var workflowRunNodesCmd = &cobra.Command{
	Use:   "nodes <workflow-id> <run-id>",
	Short: "List run nodes",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/api/workflows/"+args[0]+"/runs/"+args[1]+"/nodes", nil, nil, client.AuthAuto)
	},
}

var workflowRunEventsCmd = &cobra.Command{
	Use:   "events <workflow-id> <run-id>",
	Short: "List run events",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeID, _ := cmd.Flags().GetString("node-id")
		since, _ := cmd.Flags().GetString("since")
		limit, _ := cmd.Flags().GetInt("limit")
		query := map[string]string{
			"nodeId": nodeID,
			"since":  since,
		}
		if limit > 0 {
			query["limit"] = fmt.Sprintf("%d", limit)
		}
		return runRequest("GET", "/api/workflows/"+args[0]+"/runs/"+args[1]+"/events", query, nil, client.AuthAuto)
	},
}

func init() {
	for _, c := range []*cobra.Command{workflowCreateCmd, workflowUpdateCmd} {
		c.Flags().String("file", "", "JSON file path")
		c.Flags().String("data", "", "inline JSON body")
	}
	workflowRunStartCmd.Flags().String("trigger", "", "trigger reason")
	workflowRunEventsCmd.Flags().String("node-id", "", "filter by node id")
	workflowRunEventsCmd.Flags().String("since", "", "RFC3339 timestamp")
	workflowRunEventsCmd.Flags().Int("limit", 0, "max events")

	workflowRunCmd.AddCommand(workflowRunListCmd)
	workflowRunCmd.AddCommand(workflowRunStartCmd)
	workflowRunCmd.AddCommand(workflowRunCancelCmd)
	workflowRunCmd.AddCommand(workflowRunNodesCmd)
	workflowRunCmd.AddCommand(workflowRunEventsCmd)

	workflowCmd.AddCommand(workflowListCmd)
	workflowCmd.AddCommand(workflowGetCmd)
	workflowCmd.AddCommand(workflowCreateCmd)
	workflowCmd.AddCommand(workflowUpdateCmd)
	workflowCmd.AddCommand(workflowDeleteCmd)
	workflowCmd.AddCommand(workflowStatsCmd)
	workflowCmd.AddCommand(workflowRunCmd)
	rootCmd.AddCommand(workflowCmd)
}
