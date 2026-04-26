package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	sortField   string
	sortReverse bool
)

func init() {
	listCmd.Flags().StringVar(&sortField, "sort", "port", "Sort results by field: port, protocol, process, state")
	listCmd.Flags().BoolVar(&sortReverse, "reverse", false, "Reverse sort order")
}

// applySorting applies sorting to entries based on CLI flags.
// Returns the sorted slice or exits on invalid input.
func applySorting(cmd *cobra.Command, entries []ports.PortEntry) []ports.PortEntry {
	field, _ := cmd.Flags().GetString("sort")
	reverse, _ := cmd.Flags().GetBool("reverse")

	if field == "" {
		return entries
	}

	sorter, err := ports.NewSorter(field, reverse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	return sorter.Sort(entries)
}
