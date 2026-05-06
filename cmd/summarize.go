package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	summarizeTop  int
	summarizeJSON bool
)

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Print a concise summary of active port mappings",
	Long:  `Summarize scans active ports and displays totals by protocol, state, unique ports, and top processes.`,
	RunE:  runSummarize,
}

func init() {
	summarizeCmd.Flags().IntVar(&summarizeTop, "top", 5, "number of top processes to display")
	summarizeCmd.Flags().BoolVar(&summarizeJSON, "json", false, "output summary as JSON")
	rootCmd.AddCommand(summarizeCmd)
}

func runSummarize(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}

	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	summarizer, err := ports.NewSummarizer(summarizeTop)
	if err != nil {
		return fmt.Errorf("summarizer: %w", err)
	}

	sum := summarizer.Summarize(entries)

	if summarizeJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(sum)
	}

	fmt.Fprintf(os.Stdout, "Total entries   : %d\n", sum.Total)
	fmt.Fprintf(os.Stdout, "Unique ports    : %d\n", sum.UniquePorts)
	fmt.Fprintf(os.Stdout, "Unique processes: %d\n", sum.UniqueProcs)
	fmt.Fprintln(os.Stdout, "\nBy protocol:")
	for proto, count := range sum.ByProtocol {
		fmt.Fprintf(os.Stdout, "  %-6s %d\n", proto, count)
	}
	fmt.Fprintln(os.Stdout, "\nBy state:")
	for state, count := range sum.ByState {
		fmt.Fprintf(os.Stdout, "  %-12s %d\n", state, count)
	}
	if len(sum.TopProcesses) > 0 {
		fmt.Fprintf(os.Stdout, "\nTop %d processes:\n", summarizeTop)
		for i, p := range sum.TopProcesses {
			fmt.Fprintf(os.Stdout, "  %d. %s\n", i+1, p)
		}
	}
	return nil
}
