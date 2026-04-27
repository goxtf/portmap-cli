package ports

import (
	"fmt"
	"sort"
	"strings"
)

// SortField represents the field to sort by.
type SortField string

const (
	SortByPort     SortField = "port"
	SortByProtocol SortField = "protocol"
	SortByProcess  SortField = "process"
	SortByState    SortField = "state"
)

// ValidSortFields lists all accepted sort field values.
var ValidSortFields = []SortField{SortByPort, SortByProtocol, SortByProcess, SortByState}

// Sorter sorts port entries by a given field.
type Sorter struct {
	field   SortField
	reverse bool
}

// NewSorter creates a new Sorter. Returns an error if the field is invalid.
func NewSorter(field string, reverse bool) (*Sorter, error) {
	f := SortField(strings.ToLower(field))
	switch f {
	case SortByPort, SortByProtocol, SortByProcess, SortByState:
		// valid
	default:
		return nil, fmt.Errorf("invalid sort field %q: must be one of port, protocol, process, state", field)
	}
	return &Sorter{field: f, reverse: reverse}, nil
}

// Field returns the field this Sorter is configured to sort by.
func (s *Sorter) Field() SortField {
	return s.field
}

// Reverse returns true if the Sorter is configured to sort in descending order.
func (s *Sorter) Reverse() bool {
	return s.reverse
}

// Sort returns a sorted copy of the provided entries.
func (s *Sorter) Sort(entries []PortEntry) []PortEntry {
	result := make([]PortEntry, len(entries))
	copy(result, entries)

	sort.SliceStable(result, func(i, j int) bool {
		less := s.less(result[i], result[j])
		if s.reverse {
			return !less
		}
		return less
	})
	return result
}

func (s *Sorter) less(a, b PortEntry) bool {
	switch s.field {
	case SortByPort:
		return a.Port < b.Port
	case SortByProtocol:
		return strings.ToLower(a.Protocol) < strings.ToLower(b.Protocol)
	case SortByProcess:
		return strings.ToLower(a.Process) < strings.ToLower(b.Process)
	case SortByState:
		return strings.ToLower(a.State) < strings.ToLower(b.State)
	default:
		return a.Port < b.Port
	}
}
