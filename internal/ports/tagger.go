package ports

import "fmt"

// Tag represents a label applied to a port entry.
type Tag struct {
	Key   string
	Value string
}

// TagRule defines a condition and the tag to apply when matched.
type TagRule struct {
	Field string
	Match string
	Key   string
	Value string
}

// Tagger applies user-defined tags to port entries based on matching rules.
type Tagger struct {
	rules []TagRule
}

var validTagFields = map[string]bool{
	"protocol": true,
	"state":    true,
	"process":  true,
	"port":     true,
}

// NewTagger creates a Tagger from the given rules.
// Returns an error if any rule references an invalid field.
func NewTagger(rules []TagRule) (*Tagger, error) {
	for _, r := range rules {
		if !validTagFields[r.Field] {
			return nil, fmt.Errorf("tagger: invalid field %q", r.Field)
		}
		if r.Key == "" {
			return nil, fmt.Errorf("tagger: tag key must not be empty")
		}
	}
	return &Tagger{rules: rules}, nil
}

// Tag applies all matching rules to each entry and returns a map of entry index to tags.
func (t *Tagger) Tag(entries []PortEntry) map[int][]Tag {
	result := make(map[int][]Tag)
	for i, e := range entries {
		for _, r := range t.rules {
			if t.matches(e, r) {
				result[i] = append(result[i], Tag{Key: r.Key, Value: r.Value})
			}
		}
	}
	return result
}

func (t *Tagger) matches(e PortEntry, r TagRule) bool {
	var fieldVal string
	switch r.Field {
	case "protocol":
		fieldVal = e.Protocol
	case "state":
		fieldVal = e.State
	case "process":
		fieldVal = e.Process
	case "port":
		fieldVal = fmt.Sprintf("%d", e.Port)
	}
	return fieldVal == r.Match
}
