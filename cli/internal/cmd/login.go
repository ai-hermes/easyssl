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
	Long: `Login to an EasySSL server using email and password.
The JWT token is persisted to the local config for subsequent commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server, _ := cmd.Flags().GetString("server")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if server == "" {
			return fmt.Errorf("--server is required")
		}
		if email == "" {
			return fmt.Errorf("--email is required")
		}
		if password == "" {
			return fmt.Errorf("--password is required")
		}

		cfg.Server = server
		c := client.New(cfg)
		token, err := c.Login(email, password)
		if err != nil {
			return err
		}
		cfg.Token = token
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Println("Login successful.")
		return nil
	},
}

func init() {
	loginCmd.Flags().String("server", "", "EasySSL server URL (e.g. http://localhost:8090)")
	loginCmd.Flags().String("email", "", "user email")
	loginCmd.Flags().String("password", "", "user password")
	_ = loginCmd.MarkFlagRequired("server")
	_ = loginCmd.MarkFlagRequired("email")
	_ = loginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(loginCmd)
}
