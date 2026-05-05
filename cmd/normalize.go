package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var normalizeFields string

var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize field values of active port entries",
	Long:  `Standardize protocol (lowercase), state (uppercase), and process (trimmed) fields.`,
	RunE:  runNormalize,
}

func init() {
	normalizeCmd.Flags().StringVarP(&normalizeFields, "fields", "f", "",
		"Comma-separated fields to normalize: protocol,state,process (default: all)")
	rootCmd.AddCommand(normalizeCmd)
}

func runNormalize(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	var fields []ports.NormalizeField
	if normalizeFields != "" {
		for _, f := range strings.Split(normalizeFields, ",") {
			fields = append(fields, ports.NormalizeField(strings.TrimSpace(f)))
		}
	}

	normalizer, err := ports.NewNormalizer(fields)
	if err != nil {
		return fmt.Errorf("normalizer: %w", err)
	}

	normalized := normalizer.Normalize(entries)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(normalized); err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}
