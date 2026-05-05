package ports

import "sort"

// ScoreWeights defines the weight for each scoring dimension.
type ScoreWeights struct {
	PortBelow1024 float64
	HasProcess    float64
	TCPProtocol   float64
	ListenState   float64
}

// DefaultScoreWeights returns sensible default weights.
func DefaultScoreWeights() ScoreWeights {
	return ScoreWeights{
		PortBelow1024: 3.0,
		HasProcess:    2.0,
		TCPProtocol:   1.0,
		ListenState:   1.5,
	}
}

// ScoredEntry pairs a PortEntry with its computed score.
type ScoredEntry struct {
	Entry PortEntry
	Score float64
}

// Scorer ranks port entries by a weighted relevance score.
type Scorer struct {
	weights ScoreWeights
}

// NewScorer creates a Scorer with the given weights.
// Returns an error if all weights are zero.
func NewScorer(w ScoreWeights) (*Scorer, error) {
	if w.PortBelow1024 == 0 && w.HasProcess == 0 && w.TCPProtocol == 0 && w.ListenState == 0 {
		return nil, fmt.Errorf("scorer: at least one weight must be non-zero")
	}
	return &Scorer{weights: w}, nil
}

// Score computes a relevance score for a single PortEntry.
func (s *Scorer) Score(e PortEntry) float64 {
	var score float64
	if e.Port > 0 && e.Port < 1024 {
		score += s.weights.PortBelow1024
	}
	if e.Process != "" {
		score += s.weights.HasProcess
	}
	if e.Protocol == "tcp" {
		score += s.weights.TCPProtocol
	}
	if e.State == "LISTEN" {
		score += s.weights.ListenState
	}
	return score
}

// Rank scores all entries and returns them sorted descending by score.
func (s *Scorer) Rank(entries []PortEntry) []ScoredEntry {
	scored := make([]ScoredEntry, len(entries))
	for i, e := range entries {
		scored[i] = ScoredEntry{Entry: e, Score: s.Score(e)}
	}
	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
	return scored
}
