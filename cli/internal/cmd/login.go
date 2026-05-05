package cmd

import (
	"fmt"

	"easyssl/cli/internal/client"
	"easyssl/cli/internal/config"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the EasySSL server",
	Long: `Login to an EasySSL server using an API key.
The API key is persisted to the local config for subsequent commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server, _ := cmd.Flags().GetString("server")
		apiKey, _ := cmd.Flags().GetString("api-key")

		if server == "" {
			server = cfg.Server
		}
		if apiKey == "" {
			return fmt.Errorf("--api-key is required")
		}

		cfg.Server = server
		cfg.APIKey = apiKey
		cfg.Token = ""
		c := client.New(cfg)
		if _, err := c.Me(); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Println("Login successful.")
		return nil
	},
}

func init() {
	loginCmd.Flags().String("server", "", "EasySSL server URL (optional, defaults to https://easyssl.spotty.com.cn/)")
	loginCmd.Flags().String("api-key", "", "EasySSL API key")
	_ = loginCmd.MarkFlagRequired("api-key")
	rootCmd.AddCommand(loginCmd)
}
