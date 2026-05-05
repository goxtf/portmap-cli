package ports

import (
	"testing"
)

var rankerEntries = []PortEntry{
	{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 443, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 8080, Protocol: "udp", State: "ESTABLISHED", Process: "go"},
	{Port: 22, Protocol: "tcp", State: "LISTEN", Process: "sshd"},
}

func TestNewRankerValidFields(t *testing.T) {
	for _, field := range []string{"port", "process", "protocol", "state"} {
		_, err := NewRanker(field, false)
		if err != nil {
			t.Errorf("expected no error for field %q, got %v", field, err)
		}
	}
}

func TestNewRankerInvalidField(t *testing.T) {
	_, err := NewRanker("invalid", false)
	if err == nil {
		t.Error("expected error for invalid field, got nil")
	}
}

func TestRankByPortDescending(t *testing.T) {
	r, _ := NewRanker("port", false)
	ranked := r.Rank(rankerEntries)

	if ranked[0].Entry.Port != 8080 {
		t.Errorf("expected first rank port=8080, got %d", ranked[0].Entry.Port)
	}
	if ranked[0].Rank != 1 {
		t.Errorf("expected rank 1 for first entry, got %d", ranked[0].Rank)
	}
}

func TestRankByPortAscending(t *testing.T) {
	r, _ := NewRanker("port", true)
	ranked := r.Rank(rankerEntries)

	if ranked[0].Entry.Port != 22 {
		t.Errorf("expected first rank port=22, got %d", ranked[0].Entry.Port)
	}
}

func TestRankAssignsSequentialRanks(t *testing.T) {
	r, _ := NewRanker("port", false)
	ranked := r.Rank(rankerEntries)

	for i, re := range ranked {
		if re.Rank != i+1 {
			t.Errorf("expected rank %d at index %d, got %d", i+1, i, re.Rank)
		}
	}
}

func TestRankByProtocol(t *testing.T) {
	r, _ := NewRanker("protocol", false)
	ranked := r.Rank(rankerEntries)

	// tcp scores 2, udp scores 1; tcp entries should come first
	if ranked[0].Entry.Protocol != "tcp" {
		t.Errorf("expected tcp first, got %s", ranked[0].Entry.Protocol)
	}
}

func TestRankEmptyInput(t *testing.T) {
	r, _ := NewRanker("port", false)
	ranked := r.Rank([]PortEntry{})

	if len(ranked) != 0 {
		t.Errorf("expected empty result, got %d entries", len(ranked))
	}
}

func TestRankScoreIsPopulated(t *testing.T) {
	r, _ := NewRanker("port", false)
	ranked := r.Rank(rankerEntries)

	for _, re := range ranked {
		if re.Score != re.Entry.Port {
			t.Errorf("expected score=%d for port field, got %d", re.Entry.Port, re.Score)
		}
	}
}
