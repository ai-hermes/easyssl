package cmd

import (
	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage access credentials",
}

var accessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List accesses",
	RunE: func(cmd *cobra.Command, args []string) error {
		openapi, _ := cmd.Flags().GetBool("openapi")
		if openapi {
			return runRequest("GET", "/openapi/accesses", nil, nil, client.AuthAPIKey)
		}
		return runRequest("GET", "/api/accesses", nil, nil, client.AuthAuto)
	},
}

var accessCreateCmd = &cobra.Command{
	Use:   "create --file <json> | --data <json>",
	Short: "Create an access",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		data, _ := cmd.Flags().GetString("data")
		body, err := readJSONInput(file, data)
		if err != nil {
			return exitErr(2, err)
		}
		return runRequest("POST", "/api/accesses", nil, body, client.AuthAuto)
	},
}

var accessUpdateCmd = &cobra.Command{
	Use:   "update <id> --file <json> | --data <json>",
	Short: "Update an access",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		data, _ := cmd.Flags().GetString("data")
		body, err := readJSONInput(file, data)
		if err != nil {
			return exitErr(2, err)
		}
		return runRequest("PUT", "/api/accesses/"+args[0], nil, body, client.AuthAuto)
	},
}

var accessDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an access",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("DELETE", "/api/accesses/"+args[0], nil, nil, client.AuthAuto)
	},
}

var accessTestCmd = &cobra.Command{
	Use:   "test <id>",
	Short: "Test an access credential",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("POST", "/api/accesses/"+args[0]+"/test", nil, nil, client.AuthAuto)
	},
}

func init() {
	accessListCmd.Flags().Bool("openapi", false, "use OpenAPI endpoint /openapi/accesses")
	for _, c := range []*cobra.Command{accessCreateCmd, accessUpdateCmd} {
		c.Flags().String("file", "", "JSON file path")
		c.Flags().String("data", "", "inline JSON body")
	}

	accessCmd.AddCommand(accessListCmd)
	accessCmd.AddCommand(accessCreateCmd)
	accessCmd.AddCommand(accessUpdateCmd)
	accessCmd.AddCommand(accessDeleteCmd)
	accessCmd.AddCommand(accessTestCmd)
	rootCmd.AddCommand(accessCmd)
}
