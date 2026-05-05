package cmd

import (
	"fmt"

	"easyssl/cli/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	cfg     config.Config
)

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "easyssl",
	Short: "EasySSL CLI - manage SSL certificates from the command line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/easyssl/config.yaml)")
}
