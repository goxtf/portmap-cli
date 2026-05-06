package ports

import "fmt"

// Summary holds a human-readable digest of a set of port entries.
type Summary struct {
	Total        int            `json:"total"`
	ByProtocol   map[string]int `json:"by_protocol"`
	ByState      map[string]int `json:"by_state"`
	UniqueProcs  int            `json:"unique_processes"`
	UniquePorts  int            `json:"unique_ports"`
	TopProcesses []string       `json:"top_processes"`
}

// Summarizer produces a Summary from a slice of PortEntry values.
type Summarizer struct {
	topN int
}

// NewSummarizer returns a Summarizer that includes the top n processes.
// n must be >= 1.
func NewSummarizer(n int) (*Summarizer, error) {
	if n < 1 {
		return nil, fmt.Errorf("summarizer: n must be >= 1, got %d", n)
	}
	return &Summarizer{topN: n}, nil
}

// Summarize computes a Summary for the provided entries.
func (s *Summarizer) Summarize(entries []PortEntry) Summary {
	byProto := make(map[string]int)
	byState := make(map[string]int)
	procCount := make(map[string]int)
	portSet := make(map[int]struct{})

	for _, e := range entries {
		byProto[e.Protocol]++
		byState[e.State]++
		if e.Process != "" {
			procCount[e.Process]++
		}
		portSet[e.Port] = struct{}{}
	}

	return Summary{
		Total:        len(entries),
		ByProtocol:   byProto,
		ByState:      byState,
		UniqueProcs:  len(procCount),
		UniquePorts:  len(portSet),
		TopProcesses: topNKeys(procCount, s.topN),
	}
}

// topNKeys returns up to n keys from m ordered by descending value.
func topNKeys(m map[string]int, n int) []string {
	type kv struct {
		key string
		val int
	}
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	// simple selection sort for small slices
	for i := 0; i < len(pairs); i++ {
		max := i
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].val > pairs[max].val {
				max = j
			}
		}
		pairs[i], pairs[max] = pairs[max], pairs[i]
	}
	result := make([]string, 0, n)
	for i := 0; i < n && i < len(pairs); i++ {
		result = append(result, pairs[i].key)
	}
	return result
}
