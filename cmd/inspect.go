package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var inspectJSON bool

var inspectCmd = &cobra.Command{
	Use:   "inspect <port>",
	Short: "Deeply inspect a specific port and its bindings",
	Args:  cobra.ExactArgs(1),
	RunE:  runInspect,
}

func init() {
	inspectCmd.Flags().BoolVar(&inspectJSON, "json", false, "Output result as JSON")
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	targetPort, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid port number: %s", args[0])
	}

	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	all, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	resolver, _ := ports.NewResolver(nil)
	scorer, _ := ports.NewScorer(ports.DefaultScoreWeights())
	inspector, err := ports.NewInspector(resolver, scorer, nil, nil)
	if err != nil {
		return fmt.Errorf("inspector error: %w", err)
	}

	var results []ports.InspectionResult
	for _, e := range all {
		if e.Port == targetPort {
			results = append(results, inspector.Inspect(e, all))
		}
	}

	if len(results) == 0 {
		fmt.Fprintf(os.Stderr, "no entries found for port %d\n", targetPort)
		return nil
	}

	if inspectJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}

	for _, r := range results {
		fmt.Printf("Port:    %d\n", r.Entry.Port)
		fmt.Printf("Proto:   %s\n", r.Entry.Protocol)
		fmt.Printf("State:   %s\n", r.Entry.State)
		fmt.Printf("Process: %s\n", r.Entry.Process)
		fmt.Printf("Service: %s\n", r.ServiceName)
		fmt.Printf("Score:   %.2f\n", r.Score)
		if r.Note != "" {
			fmt.Printf("Note:    %s\n", r.Note)
		}
		fmt.Println("---")
	}
	return nil
}
