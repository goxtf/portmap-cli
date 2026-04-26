package ports_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portmap-cli/internal/ports"
)

func makeEntries() []ports.PortEntry {
	return []ports.PortEntry{
		{Protocol: ports.TCP, LocalAddress: "0.0.0.0", Port: 443, State: ports.Listen, PID: 100, Process: "caddy"},
		{Protocol: ports.UDP, LocalAddress: "127.0.0.1", Port: 5353, State: ports.Established, PID: 200, Process: "avahi"},
		{Protocol: ports.TCP, LocalAddress: "::1", Port: 6379, State: ports.Listen, PID: 300, Process: "redis"},
	}
}

func TestExporterRoundTripJSON(t *testing.T) {
	e, err := ports.NewExporter(ports.FormatJSON)
	if err != nil {
		t.Fatalf("NewExporter: %v", err)
	}
	var buf bytes.Buffer
	if err := e.Export(makeEntries(), &buf); err != nil {
		t.Fatalf("Export: %v", err)
	}
	var decoded []ports.PortEntry
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if len(decoded) != 3 {
		t.Errorf("expected 3 entries, got %d", len(decoded))
	}
	if decoded[2].Process != "redis" {
		t.Errorf("expected process 'redis', got '%s'", decoded[2].Process)
	}
}

func TestExporterTextContainsAllProcesses(t *testing.T) {
	e, _ := ports.NewExporter(ports.FormatText)
	var buf bytes.Buffer
	_ = e.Export(makeEntries(), &buf)
	out := buf.String()
	for _, proc := range []string{"caddy", "avahi", "redis"} {
		if !strings.Contains(out, proc) {
			t.Errorf("expected process '%s' in text output", proc)
		}
	}
}

func TestExporterCSVRowCount(t *testing.T) {
	e, _ := ports.NewExporter(ports.FormatCSV)
	var buf bytes.Buffer
	_ = e.Export(makeEntries(), &buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 3 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines in CSV, got %d", len(lines))
	}
}
