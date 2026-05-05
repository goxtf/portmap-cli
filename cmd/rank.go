package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	rankField   string
	rankReverse bool
	rankJSON    bool
)

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Rank active port entries by a specified field",
	Long:  `Scan active ports and rank them by port, process, protocol, or state. Outputs ranked entries with their score and rank position.`,
	RunE:  runRank,
}

func init() {
	rankCmd.Flags().StringVarP(&rankField, "field", "f", "port", "Field to rank by (port, process, protocol, state)")
	rankCmd.Flags().BoolVarP(&rankReverse, "reverse", "r", false, "Reverse ranking order (ascending score first)")
	rankCmd.Flags().BoolVar(&rankJSON, "json", false, "Output ranked entries as JSON")
	rootCmd.AddCommand(rankCmd)
}

func runRank(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("rank: failed to create scanner: %w", err)
	}

	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("rank: scan failed: %w", err)
	}

	ranker, err := ports.NewRanker(rankField, rankReverse)
	if err != nil {
		return fmt.Errorf("rank: %w", err)
	}

	ranked := ranker.Rank(entries)

	if rankJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(ranked)
	}

	fmt.Fprintf(os.Stdout, "%-6s %-8s %-10s %-20s %-12s %s\n",
		"RANK", "SCORE", "PROTOCOL", "ADDRESS", "STATE", "PROCESS")
	for _, re := range ranked {
		fmt.Fprintf(os.Stdout, "%-6d %-8d %-10s %-20s %-12s %s\n",
			re.Rank,
			re.Score,
			re.Entry.Protocol,
			fmt.Sprintf("%s:%d", re.Entry.Address, re.Entry.Port),
			re.Entry.State,
			re.Entry.Process,
		)
	}
	return nil
}
