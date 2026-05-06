package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var traceJSON bool

func init() {
	traceCmd := &cobra.Command{
		Use:   "trace <port>",
		Short: "Trace how a port entry is processed and annotated",
		Long: `Trace scans active port mappings and prints a step-by-step audit
trail showing protocol, state, owning process, resolved service name,
and any tags or labels attached to the given port.`,
		Args: cobra.ExactArgs(1),
		RunE: runTrace,
	}
	traceCmd.Flags().BoolVar(&traceJSON, "json", false, "Output trace results as JSON")
	rootCmd.AddCommand(traceCmd)
}

func runTrace(cmd *cobra.Command, args []string) error {
	var port int
	if _, err := fmt.Sscanf(args[0], "%d", &port); err != nil {
		return fmt.Errorf("invalid port argument: %s", args[0])
	}

	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner init: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	resolver, _ := ports.NewResolver(nil)
	tracer := ports.NewTracer(resolver)

	results, err := tracer.Trace(entries, port)
	if err != nil {
		return err
	}

	if traceJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	for _, r := range results {
		fmt.Fprintf(os.Stdout, "--- port %d [%s] ---\n", r.Port, r.Protocol)
		for _, h := range r.Hops {
			fmt.Fprintf(os.Stdout, "  -> %s\n", h)
		}
		for _, n := range r.Notes {
			fmt.Fprintf(os.Stdout, "  NOTE: %s\n", n)
		}
	}
	return nil
}
