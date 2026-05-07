package ports

import (
	"testing"
)

var mergerSetA = []PortEntry{
	{Protocol: "tcp", LocalPort: 80, PID: 100, Process: "nginx", State: "LISTEN"},
	{Protocol: "tcp", LocalPort: 443, PID: 101, Process: "nginx", State: "LISTEN"},
}

var mergerSetB = []PortEntry{
	{Protocol: "tcp", LocalPort: 80, PID: 100, Process: "nginx-updated", State: "LISTEN"},
	{Protocol: "udp", LocalPort: 53, PID: 200, Process: "dns", State: "UNCONN"},
}

func TestNewMergerValidStrategies(t *testing.T) {
	for _, s := range []MergeStrategy{MergeStrategyFirst, MergeStrategyLast, MergeStrategyUnion} {
		_, err := NewMerger(s)
		if err != nil {
			t.Errorf("expected no error for strategy %q, got %v", s, err)
		}
	}
}

func TestNewMergerInvalidStrategy(t *testing.T) {
	_, err := NewMerger("bogus")
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}

func TestMergeStrategyFirst(t *testing.T) {
	m, _ := NewMerger(MergeStrategyFirst)
	result := m.Merge(mergerSetA, mergerSetB)
	// port 80 should come from setA (first)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.LocalPort == 80 && e.Process != "nginx" {
			t.Errorf("first strategy: expected process nginx, got %s", e.Process)
		}
	}
}

func TestMergeStrategyLast(t *testing.T) {
	m, _ := NewMerger(MergeStrategyLast)
	result := m.Merge(mergerSetA, mergerSetB)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	for _, e := range result {
		if e.LocalPort == 80 && e.Process != "nginx-updated" {
			t.Errorf("last strategy: expected process nginx-updated, got %s", e.Process)
		}
	}
}

func TestMergeStrategyUnion(t *testing.T) {
	m, _ := NewMerger(MergeStrategyUnion)
	result := m.Merge(mergerSetA, mergerSetB)
	// union keeps all entries including duplicates
	if len(result) != len(mergerSetA)+len(mergerSetB) {
		t.Fatalf("expected %d entries, got %d", len(mergerSetA)+len(mergerSetB), len(result))
	}
}

func TestMergeEmptySets(t *testing.T) {
	m, _ := NewMerger(MergeStrategyFirst)
	result := m.Merge([]PortEntry{}, []PortEntry{})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestMergeSingleSet(t *testing.T) {
	m, _ := NewMerger(MergeStrategyFirst)
	result := m.Merge(mergerSetA)
	if len(result) != len(mergerSetA) {
		t.Fatalf("expected %d entries, got %d", len(mergerSetA), len(result))
	}
}
