package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	splitField  string
	splitValue  string
	splitJSON   bool
)

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Partition active port entries into matched and unmatched groups",
	Long: `Split scans active port mappings and partitions them into two groups:
entries where the given field matches the value (matched) and all others (unmatched).

Supported fields: protocol, state, process`,
	RunE: runSplit,
}

func init() {
	splitCmd.Flags().StringVarP(&splitField, "field", "f", "protocol", "Field to split on (protocol, state, process)")
	splitCmd.Flags().StringVarP(&splitValue, "value", "v", "", "Value to match against the chosen field")
	splitCmd.Flags().BoolVar(&splitJSON, "json", false, "Output result as JSON")
	_ = splitCmd.MarkFlagRequired("value")
	rootCmd.AddCommand(splitCmd)
}

func runSplit(cmd *cobra.Command, _ []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	splitter, err := ports.NewSplitter(splitField, splitValue)
	if err != nil {
		return err
	}

	result := splitter.Split(entries)

	if splitJSON {
		out := map[string]interface{}{
			"matched":   result.Matched,
			"unmatched": result.Unmatched,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Matched (%d):\n", len(result.Matched))
	for _, e := range result.Matched {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\t%d\t%s\t%s\n", e.Protocol, e.Port, e.State, e.Process)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\nUnmatched (%d):\n", len(result.Unmatched))
	for _, e := range result.Unmatched {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\t%d\t%s\t%s\n", e.Protocol, e.Port, e.State, e.Process)
	}
	return nil
}
