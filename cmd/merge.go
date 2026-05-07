package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var (
	mergeStrategy string
	mergeFiles    []string
	mergeJSON     bool
)

func init() {
	mergeCmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge port entry snapshots from multiple JSON files",
		Long: `Merge combines port entry lists from two or more snapshot JSON files
into a single list. Duplicate resolution is controlled by --strategy.`,
		RunE: runMerge,
	}
	mergeCmd.Flags().StringVarP(&mergeStrategy, "strategy", "s", "first",
		"Merge strategy: first, last, union")
	mergeCmd.Flags().StringArrayVarP(&mergeFiles, "file", "f", nil,
		"Input JSON snapshot files (repeat flag for multiple files)")
	mergeCmd.Flags().BoolVar(&mergeJSON, "json", false,
		"Output merged result as JSON")
	_ = mergeCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	if len(mergeFiles) < 2 {
		return fmt.Errorf("at least two --file arguments are required")
	}

	m, err := ports.NewMerger(ports.MergeStrategy(mergeStrategy))
	if err != nil {
		return err
	}

	var sets [][]ports.PortEntry
	for _, path := range mergeFiles {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		var entries []ports.PortEntry
		if err := json.Unmarshal(data, &entries); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		sets = append(sets, entries)
	}

	result := m.Merge(sets...)

	if mergeJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Merged %d entries (strategy: %s)\n", len(result), mergeStrategy)
	for _, e := range result {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\t:%d\tpid=%d\t%s\n",
			e.Protocol, e.LocalPort, e.PID, e.Process)
	}
	return nil
}
