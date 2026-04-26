// Package cmd provides the CLI commands for portmap-cli.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via ldflags.
var version = "dev"

// rootCmd is the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "portmap",
	Short: "Manage and visualize active port mappings on your machine",
	Long: `portmap-cli is a lightweight tool to inspect active port mappings
and process bindings on your local machine.

Examples:
  portmap list                     List all active port mappings
  portmap list --proto tcp         Filter by protocol
  portmap list --state LISTEN      Filter by connection state
  portmap list --process nginx     Filter by process name
  portmap export --format json     Export port mappings as JSON
  portmap export --format csv -o ports.csv`,
	Version: version,
	// Silence default error output so we can handle it ourselves.
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags available to all subcommands.
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
