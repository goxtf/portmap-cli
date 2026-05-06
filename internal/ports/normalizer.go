package ports

import (
	"fmt"
	"strings"
)

// NormalizeField defines which fields to normalize.
type NormalizeField string

const (
	NormalizeProtocol NormalizeField = "protocol"
	NormalizeState    NormalizeField = "state"
	NormalizeProcess  NormalizeField = "process"
)

var validNormalizeFields = map[NormalizeField]bool{
	NormalizeProtocol: true,
	NormalizeState:    true,
	NormalizeProcess:  true,
}

// Normalizer standardizes field values across port entries.
type Normalizer struct {
	fields map[NormalizeField]bool
}

// NewNormalizer creates a Normalizer for the given fields.
// If fields is empty, all fields are normalized.
func NewNormalizer(fields []NormalizeField) (*Normalizer, error) {
	if len(fields) == 0 {
		all := make(map[NormalizeField]bool)
		for f := range validNormalizeFields {
			all[f] = true
		}
		return &Normalizer{fields: all}, nil
	}
	selected := make(map[NormalizeField]bool)
	for _, f := range fields {
		if !validNormalizeFields[f] {
			return nil, fmt.Errorf("normalizer: invalid field %q", f)
		}
		selected[f] = true
	}
	return &Normalizer{fields: selected}, nil
}

// Normalize returns a new slice of entries with selected fields normalized.
func (n *Normalizer) Normalize(entries []PortEntry) []PortEntry {
	result := make([]PortEntry, len(entries))
	for i, e := range entries {
		if n.fields[NormalizeProtocol] {
			e.Protocol = strings.ToLower(strings.TrimSpace(e.Protocol))
		}
		if n.fields[NormalizeState] {
			e.State = strings.ToUpper(strings.TrimSpace(e.State))
		}
		if n.fields[NormalizeProcess] {
			e.Process = strings.TrimSpace(e.Process)
		}
		result[i] = e
	}
	return result
}

// Fields returns the list of fields this Normalizer is configured to normalize.
func (n *Normalizer) Fields() []NormalizeField {
	fields := make([]NormalizeField, 0, len(n.fields))
	for f := range n.fields {
		fields = append(fields, f)
	}
	return fields
}
