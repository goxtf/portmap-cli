package ports

import (
	"testing"
)

var dedupeEntries = []PortEntry{
	{Port: 80, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
	{Port: 80, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
	{Port: 443, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
	{Port: 8080, Protocol: "TCP", Process: "go", State: "LISTEN"},
	{Port: 8080, Protocol: "TCP", Process: "go", State: "LISTEN"},
	{Port: 8080, Protocol: "UDP", Process: "go", State: "LISTEN"},
}

func TestNewDeduplicatorDefaultFields(t *testing.T) {
	d, err := NewDeduplicator(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.keyFields) != 4 {
		t.Errorf("expected 4 default fields, got %d", len(d.keyFields))
	}
}

func TestNewDeduplicatorInvalidField(t *testing.T) {
	_, err := NewDeduplicator([]string{"port", "unknown"})
	if err == nil {
		t.Fatal("expected error for invalid field, got nil")
	}
}

func TestNewDeduplicatorValidFields(t *testing.T) {
	fields := []string{"port", "protocol"}
	d, err := NewDeduplicator(fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.keyFields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(d.keyFields))
	}
}

func TestDeduplicateRemovesDuplicates(t *testing.T) {
	d, _ := NewDeduplicator(nil)
	result := d.Deduplicate(dedupeEntries)

	if len(result.Unique) != 4 {
		t.Errorf("expected 4 unique entries, got %d", len(result.Unique))
	}
	if result.Removed != 2 {
		t.Errorf("expected 2 removed, got %d", result.Removed)
	}
	if len(result.Duplicates) != 2 {
		t.Errorf("expected 2 duplicates, got %d", len(result.Duplicates))
	}
}

func TestDeduplicateByPortOnly(t *testing.T) {
	d, _ := NewDeduplicator([]string{"port"})
	result := d.Deduplicate(dedupeEntries)

	// port 80 appears twice, 8080 appears three times -> 3 unique ports
	if len(result.Unique) != 3 {
		t.Errorf("expected 3 unique by port, got %d", len(result.Unique))
	}
	if result.Removed != 3 {
		t.Errorf("expected 3 removed, got %d", result.Removed)
	}
}

func TestDeduplicateEmptyInput(t *testing.T) {
	d, _ := NewDeduplicator(nil)
	result := d.Deduplicate([]PortEntry{})

	if len(result.Unique) != 0 {
		t.Errorf("expected 0 unique entries, got %d", len(result.Unique))
	}
	if result.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", result.Removed)
	}
}
