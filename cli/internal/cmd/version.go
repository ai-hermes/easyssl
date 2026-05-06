package cmd

import (
	"fmt"

	"easyssl/cli/internal/buildinfo"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI build information",
	RunE: func(cmd *cobra.Command, args []string) error {
		info := buildinfo.Get()
		if cfg.Output == "text" {
			fmt.Printf("version=%s commit=%s buildDate=%s go=%s platform=%s\n", info.Version, info.Commit, info.BuildDate, info.GoVersion, info.Platform)
			return nil
		}
		return printOutput(info)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
