package ports

import "fmt"

// SplitResult holds two partitions of entries.
type SplitResult struct {
	Matched   []PortEntry
	Unmatched []PortEntry
}

// Splitter partitions a slice of PortEntry into two groups based on a field
// predicate, similar to a stream partition operation.
type Splitter struct {
	field string
	value string
}

var validSplitFields = map[string]struct{}{
	"protocol": {},
	"state":    {},
	"process":  {},
}

// NewSplitter creates a Splitter that partitions entries where field == value.
func NewSplitter(field, value string) (*Splitter, error) {
	if _, ok := validSplitFields[field]; !ok {
		return nil, fmt.Errorf("splitter: invalid field %q: must be one of protocol, state, process", field)
	}
	if value == "" {
		return nil, fmt.Errorf("splitter: value must not be empty")
	}
	return &Splitter{field: field, value: value}, nil
}

// Split partitions entries into Matched and Unmatched groups.
func (s *Splitter) Split(entries []PortEntry) SplitResult {
	var result SplitResult
	for _, e := range entries {
		if s.matches(e) {
			result.Matched = append(result.Matched, e)
		} else {
			result.Unmatched = append(result.Unmatched, e)
		}
	}
	return result
}

func (s *Splitter) matches(e PortEntry) bool {
	switch s.field {
	case "protocol":
		return e.Protocol == s.value
	case "state":
		return e.State == s.value
	case "process":
		return e.Process == s.value
	}
	return false
}
