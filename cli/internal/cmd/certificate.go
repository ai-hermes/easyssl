package cmd

import (
	"fmt"

	"easyssl/cli/internal/client"

	"github.com/spf13/cobra"
)

var certificateCmd = &cobra.Command{
	Use:     "certificate",
	Aliases: []string{"cert"},
	Short:   "Manage certificates",
}

var certListCmd = &cobra.Command{
	Use:   "list",
	Short: "List certificates",
	RunE: func(cmd *cobra.Command, args []string) error {
		openapi, _ := cmd.Flags().GetBool("openapi")
		if openapi {
			return runRequest("GET", "/openapi/certificates", nil, nil, client.AuthAPIKey)
		}
		return runRequest("GET", "/api/certificates", nil, nil, client.AuthAuto)
	},
}

var certApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply for a certificate via OpenAPI",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		data, _ := cmd.Flags().GetString("data")
		body := map[string]any{}
		if file != "" || data != "" {
			parsed, err := readJSONInput(file, data)
			if err != nil {
				return exitErr(2, err)
			}
			body = parsed
		}
		provider, _ := cmd.Flags().GetString("provider")
		accessID, _ := cmd.Flags().GetString("access-id")
		domains, _ := cmd.Flags().GetStringSlice("domain")
		contactEmail, _ := cmd.Flags().GetString("contact-email")
		if provider != "" {
			body["provider"] = provider
		}
		if accessID != "" {
			body["accessId"] = accessID
		}
		if len(domains) > 0 {
			body["domains"] = domains
		}
		if contactEmail != "" {
			body["contactEmail"] = contactEmail
		}
		if len(body) == 0 {
			return exitErr(2, fmt.Errorf("provide request body via --file/--data or flags"))
		}
		return runRequest("POST", "/openapi/certificates/apply", nil, body, client.AuthAPIKey)
	},
}

var certDownloadCmd = &cobra.Command{
	Use:   "download <certificate-id>",
	Short: "Get certificate download info",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("cert-format")
		openapi, _ := cmd.Flags().GetBool("openapi")
		path := "/api/certificates/" + args[0] + "/download"
		auth := client.AuthAuto
		if openapi {
			path = "/openapi/certificates/" + args[0] + "/download"
			auth = client.AuthAPIKey
		}
		body := map[string]any{"format": format}
		return runRequest("POST", path, nil, body, auth)
	},
}

var certRevokeCmd = &cobra.Command{
	Use:   "revoke <certificate-id>",
	Short: "Revoke a certificate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("POST", "/api/certificates/"+args[0]+"/revoke", nil, nil, client.AuthAuto)
	},
}

var certStatusCmd = &cobra.Command{
	Use:   "status <run-id>",
	Short: "Get openapi certificate apply run status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRequest("GET", "/openapi/certificates/runs/"+args[0], nil, nil, client.AuthAPIKey)
	},
}

var certRunEventsCmd = &cobra.Command{
	Use:   "events <run-id>",
	Short: "List openapi certificate run events",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeID, _ := cmd.Flags().GetString("node-id")
		since, _ := cmd.Flags().GetString("since")
		limit, _ := cmd.Flags().GetInt("limit")
		query := map[string]string{"nodeId": nodeID, "since": since}
		if limit > 0 {
			query["limit"] = fmt.Sprintf("%d", limit)
		}
		return runRequest("GET", "/openapi/certificates/runs/"+args[0]+"/events", query, nil, client.AuthAPIKey)
	},
}

func init() {
	certListCmd.Flags().Bool("openapi", false, "use OpenAPI endpoint /openapi/certificates")
	certApplyCmd.Flags().String("file", "", "JSON file path")
	certApplyCmd.Flags().String("data", "", "inline JSON body")
	certApplyCmd.Flags().String("provider", "", "certificate provider")
	certApplyCmd.Flags().String("access-id", "", "access credential id")
	certApplyCmd.Flags().StringSlice("domain", nil, "domain list, can be repeated")
	certApplyCmd.Flags().String("contact-email", "", "contact email")
	certDownloadCmd.Flags().String("cert-format", "PEM", "download format: PEM|PFX|JKS")
	certDownloadCmd.Flags().Bool("openapi", false, "use OpenAPI endpoint")
	certRunEventsCmd.Flags().String("node-id", "", "filter by node id")
	certRunEventsCmd.Flags().String("since", "", "RFC3339 timestamp")
	certRunEventsCmd.Flags().Int("limit", 0, "max events")

	certificateCmd.AddCommand(certListCmd)
	certificateCmd.AddCommand(certApplyCmd)
	certificateCmd.AddCommand(certDownloadCmd)
	certificateCmd.AddCommand(certRevokeCmd)
	certificateCmd.AddCommand(certStatusCmd)
	certificateCmd.AddCommand(certRunEventsCmd)
	rootCmd.AddCommand(certificateCmd)
}
