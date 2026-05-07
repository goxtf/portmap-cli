package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var (
	annotateField  string
	annotateMatch  string
	annotateNote   string
	annotateJSON   bool
)

var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Attach notes to port entries matching a field/value rule",
	Long:  `Scan active ports and annotate matching entries with a custom note.`,
	RunE:  runAnnotate,
}

func init() {
	annotateCmd.Flags().StringVarP(&annotateField, "field", "f", "process", "Field to match (protocol, state, process, port)")
	annotateCmd.Flags().StringVarP(&annotateMatch, "match", "m", "", "Value to match against the field (required)")
	annotateCmd.Flags().StringVarP(&annotateNote, "note", "n", "", "Note to attach to matching entries (required)")
	annotateCmd.Flags().BoolVar(&annotateJSON, "json", false, "Output results as JSON")
	_ = annotateCmd.MarkFlagRequired("match")
	_ = annotateCmd.MarkFlagRequired("note")
	rootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner init: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	annotator, err := ports.NewAnnotator([]ports.AnnotationRule{
		{Field: annotateField, Match: annotateMatch, Note: annotateNote},
	})
	if err != nil {
		return fmt.Errorf("annotator: %w", err)
	}

	annotated := annotator.Annotate(entries)

	if annotateJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(annotated)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%-8s %-6s %-15s %-12s %s\n", "PORT", "PROTO", "STATE", "PROCESS", "NOTE")
	fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 60))
	for _, e := range annotated {
		if e.Notes != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "%-8d %-6s %-15s %-12s %s\n",
				e.Port, e.Protocol, e.State, e.Process, e.Notes)
		}
	}
	return nil
}
