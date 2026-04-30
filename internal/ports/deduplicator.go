package ports

import "fmt"

// DeduplicateResult holds the outcome of a deduplication pass.
type DeduplicateResult struct {
	Unique     []PortEntry
	Duplicates []PortEntry
	Removed    int
}

// Deduplicator removes duplicate PortEntry records based on configurable key fields.
type Deduplicator struct {
	keyFields []string
}

var validDedupeFields = map[string]bool{
	"port":     true,
	"protocol": true,
	"process":  true,
	"state":    true,
}

// NewDeduplicator creates a Deduplicator that considers the given fields when
// determining uniqueness. If no fields are provided, all four fields are used.
func NewDeduplicator(keyFields []string) (*Deduplicator, error) {
	if len(keyFields) == 0 {
		keyFields = []string{"port", "protocol", "process", "state"}
	}
	for _, f := range keyFields {
		if !validDedupeFields[f] {
			return nil, fmt.Errorf("deduplicator: invalid key field %q", f)
		}
	}
	return &Deduplicator{keyFields: keyFields}, nil
}

// Deduplicate processes entries and returns unique entries plus metadata about
// which entries were considered duplicates.
func (d *Deduplicator) Deduplicate(entries []PortEntry) DeduplicateResult {
	seen := make(map[string]bool)
	var unique []PortEntry
	var duplicates []PortEntry

	for _, e := range entries {
		key := d.buildKey(e)
		if seen[key] {
			duplicates = append(duplicates, e)
		} else {
			seen[key] = true
			unique = append(unique, e)
		}
	}

	return DeduplicateResult{
		Unique:     unique,
		Duplicates: duplicates,
		Removed:    len(duplicates),
	}
}

func (d *Deduplicator) buildKey(e PortEntry) string {
	key := ""
	for _, f := range d.keyFields {
		switch f {
		case "port":
			key += fmt.Sprintf("port=%d|", e.Port)
		case "protocol":
			key += fmt.Sprintf("proto=%s|", e.Protocol)
		case "process":
			key += fmt.Sprintf("proc=%s|", e.Process)
		case "state":
			key += fmt.Sprintf("state=%s|", e.State)
		}
	}
	return key
}
