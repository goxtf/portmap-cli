package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display aggregated statistics about active port mappings",
	Long:  `Show protocol breakdown, state distribution, unique ports, and top processes.`,
	RunE:  runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}

	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("failed to scan ports: %w", err)
	}

	collector := ports.NewStatsCollector()
	stats := collector.Collect(entries)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Total entries:\t%d\n", stats.TotalEntries)
	fmt.Fprintf(w, "Unique ports:\t%d\n", stats.UniquePorts)
	fmt.Fprintf(w, "Unique processes:\t%d\n", stats.UniqueProcesses)

	fmt.Fprintln(w, "\nBy Protocol:")
	for proto, count := range stats.ByProtocol {
		fmt.Fprintf(w, "  %s:\t%d\n", proto, count)
	}

	fmt.Fprintln(w, "\nBy State:")
	for state, count := range stats.ByState {
		fmt.Fprintf(w, "  %s:\t%d\n", state, count)
	}

	fmt.Fprintln(w, "\nTop Processes:")
	for i, pc := range stats.TopProcesses {
		fmt.Fprintf(w, "  %d. %s:\t%d ports\n", i+1, pc.Process, pc.Count)
	}

	return w.Flush()
}
