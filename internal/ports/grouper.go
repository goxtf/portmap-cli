package ports

import "sort"

// GroupBy defines the field to group entries by.
type GroupBy string

const (
	GroupByProcess  GroupBy = "process"
	GroupByProtocol GroupBy = "protocol"
	GroupByState    GroupBy = "state"
)

// GroupResult holds a group label and its associated entries.
type GroupResult struct {
	Label   string
	Entries []PortEntry
}

// Grouper groups port entries by a specified field.
type Grouper struct {
	field GroupBy
}

// NewGrouper creates a new Grouper for the given field.
// Returns an error if the field is not supported.
func NewGrouper(field string) (*Grouper, error) {
	f := GroupBy(field)
	switch f {
	case GroupByProcess, GroupByProtocol, GroupByState:
		return &Grouper{field: f}, nil
	}
	return nil, fmt.Errorf("unsupported group-by field: %q (valid: process, protocol, state)", field)
}

// Group partitions entries into named groups sorted by label.
func (g *Grouper) Group(entries []PortEntry) []GroupResult {
	buckets := make(map[string][]PortEntry)

	for _, e := range entries {
		key := g.keyFor(e)
		buckets[key] = append(buckets[key], e)
	}

	results := make([]GroupResult, 0, len(buckets))
	for label, items := range buckets {
		results = append(results, GroupResult{Label: label, Entries: items})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Label < results[j].Label
	})

	return results
}

func (g *Grouper) keyFor(e PortEntry) string {
	switch g.field {
	case GroupByProcess:
		if e.Process == "" {
			return "(unknown)"
		}
		return e.Process
	case GroupByProtocol:
		return e.Protocol
	case GroupByState:
		if e.State == "" {
			return "(none)"
		}
		return e.State
	}
	return "(unknown)"
}
