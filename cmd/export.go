package cmd

import (
	"fmt"
	"os"

	"github.com/user/portmap-cli/internal/ports"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	exportOutput string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export active port mappings to a file or stdout",
	Long:  "Scan active port mappings and export them in JSON, CSV, or text format.",
	RunE:  runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Output format: json, csv, text")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (default: stdout)")
}

func runExport(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}

	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	exporter, err := ports.NewExporter(ports.ExportFormat(exportFormat))
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
	}

	w := cmd.OutOrStdout()
	if exportOutput != "" {
		f, err := os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	if err := exporter.Export(entries, w); err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	if exportOutput != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "Exported %d entries to %s\n", len(entries), exportOutput)
	}
	return nil
}
