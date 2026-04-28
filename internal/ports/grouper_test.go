package ports

import (
	"testing"
)

func TestNewGrouperValidFields(t *testing.T) {
	for _, field := range []string{"process", "protocol", "state"} {
		g, err := NewGrouper(field)
		if err != nil {
			t.Errorf("expected no error for field %q, got %v", field, err)
		}
		if g == nil {
			t.Errorf("expected non-nil grouper for field %q", field)
		}
	}
}

func TestNewGrouperInvalidField(t *testing.T) {
	_, err := NewGrouper("invalid")
	if err == nil {
		t.Error("expected error for invalid field, got nil")
	}
}

var grouperEntries = []PortEntry{
	{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 80, State: "LISTEN", Process: "nginx"},
	{Protocol: "tcp", LocalAddress: "0.0.0.0", LocalPort: 443, State: "LISTEN", Process: "nginx"},
	{Protocol: "udp", LocalAddress: "0.0.0.0", LocalPort: 53, State: "", Process: "systemd-resolved"},
	{Protocol: "tcp", LocalAddress: "127.0.0.1", LocalPort: 5432, State: "LISTEN", Process: "postgres"},
	{Protocol: "tcp", LocalAddress: "127.0.0.1", LocalPort: 6379, State: "ESTABLISHED", Process: "redis"},
}

func TestGroupByProcess(t *testing.T) {
	g, _ := NewGrouper("process")
	results := g.Group(grouperEntries)

	if len(results) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(results))
	}
	// labels should be sorted: nginx, postgres, redis, systemd-resolved
	if results[0].Label != "nginx" {
		t.Errorf("expected first group label 'nginx', got %q", results[0].Label)
	}
	if len(results[0].Entries) != 2 {
		t.Errorf("expected 2 entries for nginx, got %d", len(results[0].Entries))
	}
}

func TestGroupByProtocol(t *testing.T) {
	g, _ := NewGrouper("protocol")
	results := g.Group(grouperEntries)

	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	if results[0].Label != "tcp" {
		t.Errorf("expected first group 'tcp', got %q", results[0].Label)
	}
	if len(results[0].Entries) != 4 {
		t.Errorf("expected 4 tcp entries, got %d", len(results[0].Entries))
	}
}

func TestGroupByState(t *testing.T) {
	g, _ := NewGrouper("state")
	results := g.Group(grouperEntries)

	labels := make(map[string]int)
	for _, r := range results {
		labels[r.Label] = len(r.Entries)
	}

	if labels["LISTEN"] != 3 {
		t.Errorf("expected 3 LISTEN entries, got %d", labels["LISTEN"])
	}
	if labels["ESTABLISHED"] != 1 {
		t.Errorf("expected 1 ESTABLISHED entry, got %d", labels["ESTABLISHED"])
	}
	if labels["(none)"] != 1 {
		t.Errorf("expected 1 (none) entry, got %d", labels["(none)"])
	}
}
