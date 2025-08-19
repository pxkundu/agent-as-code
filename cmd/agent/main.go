package main

import (
	"fmt"
	"os"

	"github.com/pxkundu/agent-as-code/internal/cmd"
)

var (
	version = "1.0.0"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	// Set version info
	cmd.SetVersionInfo(version, commit, date)

	// Execute root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
