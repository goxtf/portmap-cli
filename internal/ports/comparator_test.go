package ports

import (
	"testing"
)

var baselineEntries = []PortEntry{
	{Protocol: "tcp", Port: 80, Process: "nginx", State: "LISTEN"},
	{Protocol: "tcp", Port: 443, Process: "nginx", State: "LISTEN"},
	{Protocol: "udp", Port: 53, Process: "dns", State: "UNCONN"},
}

var currentEntries = []PortEntry{
	{Protocol: "tcp", Port: 443, Process: "nginx", State: "LISTEN"},
	{Protocol: "udp", Port: 53, Process: "dns", State: "UNCONN"},
	{Protocol: "tcp", Port: 8080, Process: "app", State: "LISTEN"},
}

func TestNewComparator(t *testing.T) {
	c := NewComparator()
	if c == nil {
		t.Fatal("expected non-nil comparator")
	}
}

func TestCompareAdded(t *testing.T) {
	c := NewComparator()
	result := c.Compare(baselineEntries, currentEntries)
	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(result.Added))
	}
	if result.Added[0].Port != 8080 {
		t.Errorf("expected added port 8080, got %d", result.Added[0].Port)
	}
}

func TestCompareRemoved(t *testing.T) {
	c := NewComparator()
	result := c.Compare(baselineEntries, currentEntries)
	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(result.Removed))
	}
	if result.Removed[0].Port != 80 {
		t.Errorf("expected removed port 80, got %d", result.Removed[0].Port)
	}
}

func TestCompareCommon(t *testing.T) {
	c := NewComparator()
	result := c.Compare(baselineEntries, currentEntries)
	if len(result.Common) != 2 {
		t.Fatalf("expected 2 common, got %d", len(result.Common))
	}
}

func TestCompareIdentical(t *testing.T) {
	c := NewComparator()
	result := c.Compare(baselineEntries, baselineEntries)
	if len(result.Added) != 0 || len(result.Removed) != 0 {
		t.Errorf("expected no diff for identical slices")
	}
	if len(result.Common) != len(baselineEntries) {
		t.Errorf("expected all entries common")
	}
}

func TestCompareEmptyBaseline(t *testing.T) {
	c := NewComparator()
	result := c.Compare([]PortEntry{}, currentEntries)
	if len(result.Added) != len(currentEntries) {
		t.Errorf("expected all current entries as added")
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removed entries")
	}
}

func TestCompareSummary(t *testing.T) {
	c := NewComparator()
	result := c.Compare(baselineEntries, currentEntries)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	expected := "Added: 1, Removed: 1, Common: 2"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
