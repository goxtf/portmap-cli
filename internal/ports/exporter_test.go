package ports

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewExporterValidFormats(t *testing.T) {
	for _, f := range []ExportFormat{FormatJSON, FormatCSV, FormatText} {
		e, err := NewExporter(f)
		if err != nil {
			t.Errorf("expected no error for format %s, got %v", f, err)
		}
		if e == nil {
			t.Errorf("expected non-nil exporter for format %s", f)
		}
	}
}

func TestNewExporterInvalidFormat(t *testing.T) {
	_, err := NewExporter("xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func exporterSampleEntries() []PortEntry {
	return []PortEntry{
		{Protocol: TCP, LocalAddress: "127.0.0.1", Port: 8080, State: Listen, PID: 1234, Process: "nginx"},
		{Protocol: UDP, LocalAddress: "0.0.0.0", Port: 53, State: Established, PID: 5678, Process: "dnsmasq"},
	}
}

func TestExportJSON(t *testing.T) {
	e, _ := NewExporter(FormatJSON)
	var buf bytes.Buffer
	if err := e.Export(exporterSampleEntries(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []PortEntry
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestExportCSV(t *testing.T) {
	e, _ := NewExporter(FormatCSV)
	var buf bytes.Buffer
	if err := e.Export(exporterSampleEntries(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid CSV output: %v", err)
	}
	// header + 2 data rows
	if len(records) != 3 {
		t.Errorf("expected 3 records (header+2), got %d", len(records))
	}
	if records[0][0] != "Protocol" {
		t.Errorf("expected header 'Protocol', got %s", records[0][0])
	}
}

func TestExportText(t *testing.T) {
	e, _ := NewExporter(FormatText)
	var buf bytes.Buffer
	if err := e.Export(exporterSampleEntries(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "nginx") {
		t.Error("expected 'nginx' in text output")
	}
	if !strings.Contains(out, "dnsmasq") {
		t.Error("expected 'dnsmasq' in text output")
	}
}
