package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var (
	listProtocol string
	listState    string
	listProcess  string
	listPort     int
	listFormat   string
)

// listCmd represents the list command which scans and displays active port mappings.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List active port mappings on this machine",
	Long: `Scan and display all active port mappings and their associated processes.

Results can be filtered by protocol, state, process name, or port number.
Output can be formatted as a table (default), JSON, or CSV.

Examples:
  portmap-cli list
  portmap-cli list --protocol tcp
  portmap-cli list --state LISTEN
  portmap-cli list --process nginx
  portmap-cli list --port 8080
  portmap-cli list --format json`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listProtocol, "protocol", "p", "", "Filter by protocol (tcp, udp)")
	listCmd.Flags().StringVarP(&listState, "state", "s", "", "Filter by connection state (e.g. LISTEN, ESTABLISHED)")
	listCmd.Flags().StringVarP(&listProcess, "process", "n", "", "Filter by process name")
	listCmd.Flags().IntVarP(&listPort, "port", "P", 0, "Filter by port number")
	listCmd.Flags().StringVarP(&listFormat, "format", "f", "table", "Output format: table, json, csv")
}

// runList executes the list command: scans ports, applies filters, and renders output.
func runList(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("failed to initialise scanner: %w", err)
	}

	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	filter := ports.NewFilter(entries)

	if listProtocol != "" {
		filter = filter.ByProtocol(listProtocol)
	}
	if listState != "" {
		filter = filter.ByState(listState)
	}
	if listProcess != "" {
		filter = filter.ByProcess(listProcess)
	}
	if listPort > 0 {
		filter = filter.ByPort(listPort)
	}

	filtered := filter.Results()

	if len(filtered) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No matching port mappings found.")
		return nil
	}

	switch listFormat {
	case "table":
		formatter := ports.NewFormatter(cmd.OutOrStdout())
		formatter.RenderTable(filtered)
	case "json", "csv":
		exporter, err := ports.NewExporter(listFormat)
		if err != nil {
			return err
		}
		if err := exporter.Export(filtered, os.Stdout); err != nil {
			return fmt.Errorf("export failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown format %q: choose table, json, or csv", listFormat)
	}

	return nil
}
