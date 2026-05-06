package ports

import (
	"testing"
	"time"
)

var prunerSampleEntries = []PortEntry{
	{Port: 80, Protocol: "TCP", State: "LISTEN", Process: "nginx"},
	{Port: 443, Protocol: "TCP", State: "LISTEN", Process: "nginx"},
	{Port: 9000, Protocol: "TCP", State: "CLOSED", Process: "php-fpm"},
	{Port: 3306, Protocol: "TCP", State: "CLOSE_WAIT", Process: ""},
	{Port: 5432, Protocol: "TCP", State: "LISTEN", Process: ""},
}

func TestNewPrunerNoCriteria(t *testing.T) {
	_, err := NewPruner(PruneOptions{})
	if err == nil {
		t.Fatal("expected error when no criteria set")
	}
}

func TestNewPrunerValid(t *testing.T) {
	p, err := NewPruner(PruneOptions{RemoveClosed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Pruner")
	}
}

func TestPruneRemovesClosed(t *testing.T) {
	p, _ := NewPruner(PruneOptions{RemoveClosed: true})
	out := p.Apply(prunerSampleEntries)
	for _, e := range out {
		if isClosed(e.State) {
			t.Errorf("entry with state %q should have been pruned", e.State)
		}
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestPruneRemovesNoProcess(t *testing.T) {
	p, _ := NewPruner(PruneOptions{RemoveNoProcess: true})
	out := p.Apply(prunerSampleEntries)
	for _, e := range out {
		if e.Process == "" {
			t.Errorf("entry with no process should have been pruned (port %d)", e.Port)
		}
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestPruneOlderThan(t *testing.T) {
	now := time.Now()
	entries := []PortEntry{
		{Port: 80, Process: "nginx", LastSeen: now.Add(-10 * time.Minute)},
		{Port: 443, Process: "nginx", LastSeen: now.Add(-2 * time.Minute)},
		{Port: 8080, Process: "app", LastSeen: now.Add(-30 * time.Second)},
	}
	p, _ := NewPruner(PruneOptions{
		OlderThan: 5 * time.Minute,
		now:       func() time.Time { return now },
	})
	out := p.Apply(entries)
	if len(out) != 2 {
		t.Errorf("expected 2 entries after pruning, got %d", len(out))
	}
}

func TestPruneCombinedCriteria(t *testing.T) {
	p, _ := NewPruner(PruneOptions{RemoveClosed: true, RemoveNoProcess: true})
	out := p.Apply(prunerSampleEntries)
	// port 9000 (CLOSED), port 3306 (CLOSE_WAIT + no process), port 5432 (no process)
	if len(out) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out))
	}
}
