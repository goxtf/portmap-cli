package ports

import "fmt"

// MergeStrategy defines how duplicate entries are resolved during a merge.
type MergeStrategy string

const (
	MergeStrategyFirst MergeStrategy = "first"
	MergeStrategyLast  MergeStrategy = "last"
	MergeStrategyUnion MergeStrategy = "union"
)

var validMergeStrategies = map[MergeStrategy]bool{
	MergeStrategyFirst: true,
	MergeStrategyLast:  true,
	MergeStrategyUnion: true,
}

// Merger combines multiple slices of PortEntry into one, resolving
// duplicates according to the configured strategy.
type Merger struct {
	strategy MergeStrategy
	keyFields []string
}

// NewMerger creates a Merger with the given strategy.
// Valid strategies: "first", "last", "union".
func NewMerger(strategy MergeStrategy) (*Merger, error) {
	if !validMergeStrategies[strategy] {
		return nil, fmt.Errorf("invalid merge strategy %q: must be one of first, last, union", strategy)
	}
	return &Merger{
		strategy:  strategy,
		keyFields: []string{"proto", "local_port", "pid"},
	}, nil
}

// Merge combines the provided entry slices into a single deduplicated slice.
func (m *Merger) Merge(sets ...[]PortEntry) []PortEntry {
	if m.strategy == MergeStrategyUnion {
		return m.mergeUnion(sets...)
	}
	seen := make(map[string]bool)
	var result []PortEntry
	for _, entries := range sets {
		for _, e := range entries {
			k := mergeKey(e)
			if m.strategy == MergeStrategyFirst {
				if !seen[k] {
					seen[k] = true
					result = append(result, e)
				}
			} else { // last
				seen[k] = true
				result = upsert(result, k, e)
			}
		}
	}
	return result
}

func (m *Merger) mergeUnion(sets ...[]PortEntry) []PortEntry {
	var all []PortEntry
	for _, entries := range sets {
		all = append(all, entries...)
	}
	return all
}

func mergeKey(e PortEntry) string {
	return fmt.Sprintf("%s|%d|%d", e.Protocol, e.LocalPort, e.PID)
}

func upsert(entries []PortEntry, key string, e PortEntry) []PortEntry {
	for i, ex := range entries {
		if mergeKey(ex) == key {
			entries[i] = e
			return entries
		}
	}
	return append(entries, e)
}
