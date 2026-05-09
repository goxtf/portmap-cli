package ports

import (
	"fmt"
	"math/rand"
)

// Sampler randomly samples a subset of port entries.
type Sampler struct {
	n    int
	seed int64
	mode string // "random" or "first" or "last"
}

// NewSampler creates a Sampler that selects up to n entries.
// mode must be one of: "random", "first", "last".
// seed is used for reproducible random sampling.
func NewSampler(n int, mode string, seed int64) (*Sampler, error) {
	if n <= 0 {
		return nil, fmt.Errorf("sampler: n must be greater than 0, got %d", n)
	}
	switch mode {
	case "random", "first", "last":
		// valid
	default:
		return nil, fmt.Errorf("sampler: unsupported mode %q; must be one of: random, first, last", mode)
	}
	return &Sampler{n: n, mode: mode, seed: seed}, nil
}

// Sample returns up to n entries from the input slice according to the configured mode.
func (s *Sampler) Sample(entries []PortEntry) []PortEntry {
	if len(entries) == 0 {
		return []PortEntry{}
	}
	switch s.mode {
	case "first":
		return s.sampleFirst(entries)
	case "last":
		return s.sampleLast(entries)
	default:
		return s.sampleRandom(entries)
	}
}

func (s *Sampler) sampleFirst(entries []PortEntry) []PortEntry {
	if s.n >= len(entries) {
		return entries
	}
	return entries[:s.n]
}

func (s *Sampler) sampleLast(entries []PortEntry) []PortEntry {
	if s.n >= len(entries) {
		return entries
	}
	return entries[len(entries)-s.n:]
}

func (s *Sampler) sampleRandom(entries []PortEntry) []PortEntry {
	if s.n >= len(entries) {
		return entries
	}
	//nolint:gosec
	r := rand.New(rand.NewSource(s.seed))
	indices := r.Perm(len(entries))
	result := make([]PortEntry, s.n)
	for i := 0; i < s.n; i++ {
		result[i] = entries[indices[i]]
	}
	return result
}
