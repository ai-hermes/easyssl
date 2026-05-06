package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"easyssl/cli/internal/client"
	"easyssl/cli/internal/config"

	"github.com/spf13/cobra"
)

var (
	cfg config.Config

	globalServer  string
	globalAPIKey  string
	globalToken   string
	globalOutput  string
	globalVerbose bool
	globalTrace   bool
	globalTimeout int
)

// ExitError encodes an explicit process exit code.
type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func exitErr(code int, err error) error {
	if err == nil {
		return nil
	}
	return ExitError{Code: code, Err: err}
}

var rootCmd = &cobra.Command{
	Use:   "easyssl",
	Short: "EasySSL CLI - manage SSL certificates from the command line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return exitErr(5, fmt.Errorf("load config: %w", err))
		}
		if globalServer != "" {
			cfg.Server = globalServer
		}
		if globalAPIKey != "" {
			cfg.APIKey = globalAPIKey
		}
		if globalToken != "" {
			cfg.Token = globalToken
		}
		if globalOutput != "" {
			cfg.Output = globalOutput
		}
		if globalTimeout > 0 {
			cfg.Timeout = globalTimeout
		}
		if globalTrace {
			cfg.Trace = true
		}
		if cfg.Output == "" {
			cfg.Output = "json"
		}
		if !isOutputFormatSupported(cfg.Output) {
			return exitErr(2, fmt.Errorf("unsupported --format %q", cfg.Output))
		}
		return nil
	},
}

// Execute adds all child commands to the root command and runs it.
func Execute() error {
	err := rootCmd.Execute()
	if err == nil {
		return nil
	}
	var ee ExitError
	if ok := AsExitError(err, &ee); ok {
		fmt.Fprintf(os.Stderr, "error: %v\n", ee.Err)
		os.Exit(ee.Code)
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(5)
	return err
}

func AsExitError(err error, target *ExitError) bool {
	ee, ok := err.(ExitError)
	if !ok {
		return false
	}
	*target = ee
	return true
}

func requireAuth() error {
	if strings.TrimSpace(cfg.APIKey) == "" && strings.TrimSpace(cfg.Token) == "" {
		return exitErr(3, fmt.Errorf("missing credentials; run 'easyssl login --api-key <key>' or set --api-key/--token"))
	}
	return nil
}

func newAPIClient() (*client.Client, error) {
	c, err := client.New(cfg, client.Options{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
		Verbose: globalVerbose,
		Trace:   globalTrace || cfg.Trace,
		Stderr:  os.Stderr,
	})
	if err != nil {
		return nil, exitErr(2, err)
	}
	return c, nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&globalServer, "server", "", "EasySSL server URL")
	rootCmd.PersistentFlags().StringVar(&globalAPIKey, "api-key", "", "API key override")
	rootCmd.PersistentFlags().StringVar(&globalToken, "token", "", "Bearer token override")
	rootCmd.PersistentFlags().StringVar(&globalOutput, "format", "", "output format: json|text")
	rootCmd.PersistentFlags().BoolVar(&globalVerbose, "verbose", false, "enable verbose diagnostics")
	rootCmd.PersistentFlags().BoolVar(&globalTrace, "trace", false, "include response body snippets in verbose logs")
	rootCmd.PersistentFlags().IntVar(&globalTimeout, "timeout", 0, "HTTP timeout in seconds")
}
