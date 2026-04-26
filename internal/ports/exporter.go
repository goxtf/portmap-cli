package ports

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ExportFormat defines the supported export formats.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

// Exporter handles serializing port entries to various formats.
type Exporter struct {
	format ExportFormat
}

// NewExporter creates a new Exporter for the given format.
func NewExporter(format ExportFormat) (*Exporter, error) {
	switch format {
	case FormatJSON, FormatCSV, FormatText:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// Export writes entries to w in the configured format.
func (e *Exporter) Export(entries []PortEntry, w io.Writer) error {
	switch e.format {
	case FormatJSON:
		return exportJSON(entries, w)
	case FormatCSV:
		return exportCSV(entries, w)
	case FormatText:
		return exportText(entries, w)
	default:
		return fmt.Errorf("unsupported format: %s", e.format)
	}
}

func exportJSON(entries []PortEntry, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func exportCSV(entries []PortEntry, w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"Protocol", "LocalAddress", "Port", "State", "PID", "Process"}); err != nil {
		return err
	}
	for _, e := range entries {
		row := []string{
			string(e.Protocol),
			e.LocalAddress,
			fmt.Sprintf("%d", e.Port),
			string(e.State),
			fmt.Sprintf("%d", e.PID),
			e.Process,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportText(entries []PortEntry, w io.Writer) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-6s %-20s %-6s %-12s %-8s %s\n",
		"PROTO", "ADDRESS", "PORT", "STATE", "PID", "PROCESS"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-6s %-20s %-6d %-12s %-8d %s\n",
			e.Protocol, e.LocalAddress, e.Port, e.State, e.PID, e.Process))
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}
