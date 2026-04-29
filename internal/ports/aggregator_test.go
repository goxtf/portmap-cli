package ports

import (
	"testing"
)

var aggregatorEntries = []PortEntry{
	{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 443, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 8080, Protocol: "tcp", State: "ESTABLISHED", Process: "go"},
	{Port: 5432, Protocol: "tcp", State: "LISTEN", Process: "postgres"},
	{Port: 53, Protocol: "udp", State: "LISTEN", Process: "dns"},
	{Port: 9000, Protocol: "tcp", State: "ESTABLISHED", Process: "go"},
}

func TestNewAggregatorValidField(t *testing.T) {
	for _, field := range []string{"process", "protocol", "state"} {
		a, err := NewAggregator(field)
		if err != nil {
			t.Errorf("expected no error for field %q, got %v", field, err)
		}
		if a == nil {
			t.Errorf("expected non-nil aggregator for field %q", field)
		}
	}
}

func TestNewAggregatorInvalidField(t *testing.T) {
	_, err := NewAggregator("unknown")
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestAggregateByProcess(t *testing.T) {
	a, _ := NewAggregator("process")
	results := a.Aggregate(aggregatorEntries)

	if len(results) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(results))
	}
	// "go" has 2 entries — should be first
	if results[0].Key != "go" {
		t.Errorf("expected first key to be 'go', got %q", results[0].Key)
	}
	if results[0].Count != 2 {
		t.Errorf("expected count 2 for 'go', got %d", results[0].Count)
	}
}

func TestAggregateByProtocol(t *testing.T) {
	a, _ := NewAggregator("protocol")
	results := a.Aggregate(aggregatorEntries)

	if len(results) != 2 {
		t.Fatalf("expected 2 protocol groups, got %d", len(results))
	}
	if results[0].Key != "tcp" {
		t.Errorf("expected 'tcp' first, got %q", results[0].Key)
	}
	if results[0].Count != 5 {
		t.Errorf("expected count 5 for tcp, got %d", results[0].Count)
	}
}

func TestAggregateByState(t *testing.T) {
	a, _ := NewAggregator("state")
	results := a.Aggregate(aggregatorEntries)

	if len(results) != 2 {
		t.Fatalf("expected 2 state groups, got %d", len(results))
	}
}

func TestAggregatePortsAreSorted(t *testing.T) {
	a, _ := NewAggregator("process")
	results := a.Aggregate(aggregatorEntries)

	for _, r := range results {
		for i := 1; i < len(r.Ports); i++ {
			if r.Ports[i] < r.Ports[i-1] {
				t.Errorf("ports not sorted for key %q", r.Key)
			}
		}
	}
}

func TestAggregateEmptyEntries(t *testing.T) {
	a, _ := NewAggregator("process")
	results := a.Aggregate([]PortEntry{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
