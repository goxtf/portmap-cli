package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	enrichGroupBy string
)

var enrichCmd = &cobra.Command{
	Use:   "enrich",
	Short: "Display port entries enriched with service names and optional grouping",
	Long:  `Resolves well-known service names for each port and optionally groups entries by process, protocol, or state.`,
	RunE:  runEnrich,
}

func init() {
	enrichCmd.Flags().StringVarP(&enrichGroupBy, "group", "g", "", "Group entries by field: process, protocol, state")
	rootCmd.AddCommand(enrichCmd)
}

func runEnrich(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("enrich: failed to create scanner: %w", err)
	}

	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("enrich: scan failed: %w", err)
	}

	resolver, err := ports.NewResolver(nil)
	if err != nil {
		return fmt.Errorf("enrich: failed to create resolver: %w", err)
	}

	enricher, err := ports.NewEnricher(resolver, enrichGroupBy)
	if err != nil {
		return fmt.Errorf("enrich: %w", err)
	}

	enriched := enricher.Enrich(entries)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PROTO\tLOCAL PORT\tPID\tPROCESS\tSTATE\tSERVICE\tGROUP")
	for _, e := range enriched {
		fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\t%s\t%s\n",
			e.Protocol, e.LocalPort, e.PID, e.Process, e.State, e.ServiceName, e.Group)
	}
	return w.Flush()
}
