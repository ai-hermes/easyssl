package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func isOutputFormatSupported(s string) bool {
	switch s {
	case "json", "text":
		return true
	default:
		return false
	}
}

func printOutput(v any) error {
	switch cfg.Output {
	case "text":
		fmt.Printf("%v\n", v)
		return nil
	case "json", "":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	default:
		return fmt.Errorf("unsupported output format: %s", cfg.Output)
	}
}

func parseAPIError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case containsAny(msg, []string{"status=401", "code=401", "status=403", "code=403"}):
		return exitErr(3, err)
	case containsAny(msg, []string{"status=404", "code=404"}):
		return exitErr(4, err)
	case containsAny(msg, []string{"invalid", "required", "unsupported", "parse", "marshal"}):
		return exitErr(2, err)
	default:
		return exitErr(5, err)
	}
}

func containsAny(s string, xs []string) bool {
	for _, x := range xs {
		if x != "" && strings.Contains(s, x) {
			return true
		}
	}
	return false
}
