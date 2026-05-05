package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/portmap-cli/internal/ports"
)

var (
	pipelineProtocol string
	pipelineState    string
	pipelineProcess  string
	pipelineSortBy   string
	pipelineReverse  bool
	pipelineDedup    bool
)

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Run a configurable processing pipeline over active port mappings",
	Long: `pipeline scans active port mappings and runs them through an ordered
chain of processing steps: filter → deduplicate → sort → display.`,
	RunE: runPipeline,
}

func init() {
	pipelineCmd.Flags().StringVarP(&pipelineProtocol, "protocol", "p", "", "filter by protocol (tcp/udp)")
	pipelineCmd.Flags().StringVarP(&pipelineState, "state", "s", "", "filter by state")
	pipelineCmd.Flags().StringVar(&pipelineProcess, "process", "", "filter by process name")
	pipelineCmd.Flags().StringVar(&pipelineSortBy, "sort", "port", "sort field (port|process|protocol|state)")
	pipelineCmd.Flags().BoolVarP(&pipelineReverse, "reverse", "r", false, "reverse sort order")
	pipelineCmd.Flags().BoolVar(&pipelineDedup, "dedup", false, "remove duplicate entries")
	rootCmd.AddCommand(pipelineCmd)
}

func runPipeline(cmd *cobra.Command, _ []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	var steps []ports.PipelineStep

	// Step 1: filter
	if pipelineProtocol != "" || pipelineState != "" || pipelineProcess != "" {
		f, ferr := ports.NewFilter(pipelineProtocol, pipelineState, pipelineProcess)
		if ferr != nil {
			return fmt.Errorf("filter: %w", ferr)
		}
		steps = append(steps, func(e []ports.PortEntry) ([]ports.PortEntry, error) {
			return f.Apply(e), nil
		})
	}

	// Step 2: deduplicate
	if pipelineDedup {
		d, derr := ports.NewDeduplicator([]string{"port", "protocol", "process"})
		if derr != nil {
			return fmt.Errorf("deduplicator: %w", derr)
		}
		steps = append(steps, func(e []ports.PortEntry) ([]ports.PortEntry, error) {
			return d.Deduplicate(e), nil
		})
	}

	// Step 3: sort
	s, serr := ports.NewSorter(pipelineSortBy, pipelineReverse)
	if serr != nil {
		return fmt.Errorf("sorter: %w", serr)
	}
	steps = append(steps, func(e []ports.PortEntry) ([]ports.PortEntry, error) {
		return s.Sort(e), nil
	})

	pipeline, perr := ports.NewPipeline(steps...)
	if perr != nil {
		return fmt.Errorf("pipeline: %w", perr)
	}

	result, rerr := pipeline.Run(entries)
	if rerr != nil {
		return fmt.Errorf("pipeline run: %w", rerr)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%-8s %-8s %-20s %-10s\n", "PROTO", "PORT", "PROCESS", "STATE")
	for _, e := range result {
		fmt.Fprintf(cmd.OutOrStdout(), "%-8s %-8d %-20s %-10s\n", e.Protocol, e.Port, e.Process, e.State)
	}
	if len(result) == 0 {
		fmt.Fprintln(os.Stderr, "no entries matched")
	}
	return nil
}
