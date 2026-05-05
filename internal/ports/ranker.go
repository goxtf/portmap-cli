package ports

import (
	"fmt"
	"sort"
)

// RankedEntry pairs a PortEntry with its computed rank score.
type RankedEntry struct {
	Entry PortEntry
	Score int
	Rank  int
}

// Ranker assigns a rank to each entry based on a scoring field.
type Ranker struct {
	field   string
	reverse bool
}

var validRankFields = map[string]bool{
	"port":     true,
	"process":  true,
	"protocol": true,
	"state":    true,
}

// NewRanker creates a Ranker that ranks entries by the given field.
// If reverse is true, lower-scoring entries receive a higher rank.
func NewRanker(field string, reverse bool) (*Ranker, error) {
	if !validRankFields[field] {
		return nil, fmt.Errorf("ranker: unsupported field %q; valid fields: port, process, protocol, state", field)
	}
	return &Ranker{field: field, reverse: reverse}, nil
}

// Rank scores and ranks a slice of PortEntry values.
func (r *Ranker) Rank(entries []PortEntry) []RankedEntry {
	type scored struct {
		entry PortEntry
		score int
	}

	scored_entries := make([]scored, len(entries))
	for i, e := range entries {
		scored_entries[i] = scored{entry: e, score: r.score(e)}
	}

	sort.SliceStable(scored_entries, func(i, j int) bool {
		if r.reverse {
			return scored_entries[i].score < scored_entries[j].score
		}
		return scored_entries[i].score > scored_entries[j].score
	})

	result := make([]RankedEntry, len(scored_entries))
	for i, s := range scored_entries {
		result[i] = RankedEntry{
			Entry: s.entry,
			Score: s.score,
			Rank:  i + 1,
		}
	}
	return result
}

// score returns a numeric score for an entry based on the configured field.
func (r *Ranker) score(e PortEntry) int {
	switch r.field {
	case "port":
		return e.Port
	case "process":
		return len(e.Process)
	case "protocol":
		if e.Protocol == "tcp" {
			return 2
		}
		return 1
	case "state":
		if e.State == "LISTEN" {
			return 2
		}
		return 1
	default:
		return 0
	}
}
