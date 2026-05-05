package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

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
		if cfg.Server == "" {
			return fmt.Errorf("not logged in; run 'easyssl login' first")
		}
		c := client.New(cfg)
		data, err := c.ListCertificates()
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

var certApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply for a new certificate via OpenAPI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Server == "" {
			return fmt.Errorf("not configured; run 'easyssl login' or set --api-key")
		}
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			apiKey = cfg.APIKey
		}
		if apiKey == "" {
			return fmt.Errorf("--api-key is required for OpenAPI endpoints")
		}

		workflowID, _ := cmd.Flags().GetString("workflow")
		if workflowID == "" {
			return fmt.Errorf("--workflow is required")
		}

		c := client.New(cfg)
		c.SetAPIKey(apiKey)
		data, err := c.ApplyCertificate(map[string]any{
			"workflowId": workflowID,
		})
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

var certDownloadCmd = &cobra.Command{
	Use:   "download <certificate-id>",
	Short: "Download a certificate archive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Server == "" {
			return fmt.Errorf("not configured; run 'easyssl login' or set --api-key")
		}
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			apiKey = cfg.APIKey
		}
		if apiKey == "" {
			return fmt.Errorf("--api-key is required for OpenAPI endpoints")
		}

		out, _ := cmd.Flags().GetString("output")
		c := client.New(cfg)
		c.SetAPIKey(apiKey)
		data, err := c.DownloadCertificate(args[0])
		if err != nil {
			return err
		}
		if out == "" {
			out = args[0] + ".zip"
		}
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		fmt.Printf("Downloaded to %s\n", out)
		return nil
	},
}

var certStatusCmd = &cobra.Command{
	Use:   "status <run-id>",
	Short: "Get certificate run status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.Server == "" {
			return fmt.Errorf("not configured; run 'easyssl login' or set --api-key")
		}
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			apiKey = cfg.APIKey
		}
		if apiKey == "" {
			return fmt.Errorf("--api-key is required for OpenAPI endpoints")
		}

		c := client.New(cfg)
		c.SetAPIKey(apiKey)
		data, err := c.GetCertificateRun(args[0])
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
	certApplyCmd.Flags().String("workflow", "", "workflow ID to use for certificate application")
	certApplyCmd.Flags().String("api-key", "", "OpenAPI key (or set in config)")
	certDownloadCmd.Flags().String("output", "", "output file path")
	certDownloadCmd.Flags().String("api-key", "", "OpenAPI key (or set in config)")
	certStatusCmd.Flags().String("api-key", "", "OpenAPI key (or set in config)")

	certificateCmd.AddCommand(certListCmd)
	certificateCmd.AddCommand(certApplyCmd)
	certificateCmd.AddCommand(certDownloadCmd)
	certificateCmd.AddCommand(certStatusCmd)
	rootCmd.AddCommand(certificateCmd)
}
