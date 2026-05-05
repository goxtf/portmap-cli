package ports

import (
	"testing"
)

func TestNewScorerAllZeroWeights(t *testing.T) {
	_, err := NewScorer(ScoreWeights{})
	if err == nil {
		t.Fatal("expected error for all-zero weights")
	}
}

func TestNewScorerValid(t *testing.T) {
	s, err := NewScorer(DefaultScoreWeights())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil scorer")
	}
}

func TestScorePrivilegedPort(t *testing.T) {
	s, _ := NewScorer(ScoreWeights{PortBelow1024: 3.0})
	e := PortEntry{Port: 80}
	if got := s.Score(e); got != 3.0 {
		t.Errorf("expected 3.0, got %f", got)
	}
}

func TestScoreHighPort(t *testing.T) {
	s, _ := NewScorer(ScoreWeights{PortBelow1024: 3.0})
	e := PortEntry{Port: 8080}
	if got := s.Score(e); got != 0.0 {
		t.Errorf("expected 0.0, got %f", got)
	}
}

func TestScoreHasProcess(t *testing.T) {
	s, _ := NewScorer(ScoreWeights{HasProcess: 2.0})
	e := PortEntry{Process: "nginx"}
	if got := s.Score(e); got != 2.0 {
		t.Errorf("expected 2.0, got %f", got)
	}
}

func TestScoreTCPProtocol(t *testing.T) {
	s, _ := NewScorer(ScoreWeights{TCPProtocol: 1.0})
	e := PortEntry{Protocol: "tcp"}
	if got := s.Score(e); got != 1.0 {
		t.Errorf("expected 1.0, got %f", got)
	}
}

func TestScoreListenState(t *testing.T) {
	s, _ := NewScorer(ScoreWeights{ListenState: 1.5})
	e := PortEntry{State: "LISTEN"}
	if got := s.Score(e); got != 1.5 {
		t.Errorf("expected 1.5, got %f", got)
	}
}

func TestRankOrderDescending(t *testing.T) {
	s, _ := NewScorer(DefaultScoreWeights())
	entries := []PortEntry{
		{Port: 8080, Protocol: "udp", State: "", Process: ""},
		{Port: 22, Protocol: "tcp", State: "LISTEN", Process: "sshd"},
		{Port: 443, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	}
	ranked := s.Rank(entries)
	if len(ranked) != 3 {
		t.Fatalf("expected 3 results, got %d", len(ranked))
	}
	for i := 1; i < len(ranked); i++ {
		if ranked[i].Score > ranked[i-1].Score {
			t.Errorf("results not sorted descending at index %d", i)
		}
	}
}

func TestRankEmptyInput(t *testing.T) {
	s, _ := NewScorer(DefaultScoreWeights())
	ranked := s.Rank([]PortEntry{})
	if len(ranked) != 0 {
		t.Errorf("expected empty result, got %d entries", len(ranked))
	}
}
