package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var aggregateField string

var aggregateCmd = &cobra.Command{
	Use:   "aggregate",
	Short: "Aggregate port entries by process, protocol, or state",
	Long:  `Group and count active port mappings by a chosen field, showing port and protocol distribution per group.`,
	RunE:  runAggregate,
}

func init() {
	aggregateCmd.Flags().StringVarP(&aggregateField, "by", "b", "process", "Field to aggregate by (process, protocol, state)")
	rootCmd.AddCommand(aggregateCmd)
}

func runAggregate(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}

	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	agg, err := ports.NewAggregator(aggregateField)
	if err != nil {
		return fmt.Errorf("invalid aggregation field %q: use process, protocol, or state", aggregateField)
	}

	results := agg.Aggregate(entries)
	if len(results) == 0 {
		fmt.Fprintln(os.Stdout, "No entries found.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "%-20s  %6s  %s\n", strings.ToUpper(aggregateField), "COUNT", "PORTS")
	fmt.Fprintln(os.Stdout, strings.Repeat("-", 60))
	for _, r := range results {
		portStrs := make([]string, len(r.Ports))
		for i, p := range r.Ports {
			portStrs[i] = fmt.Sprintf("%d", p)
		}
		portList := strings.Join(portStrs, ", ")
		if len(portList) > 30 {
			portList = portList[:27] + "..."
		}
		fmt.Fprintf(os.Stdout, "%-20s  %6d  %s\n", r.Key, r.Count, portList)
	}
	return nil
}
