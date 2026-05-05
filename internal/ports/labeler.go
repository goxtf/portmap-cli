package ports

import "fmt"

// LabelRule defines a condition and label to apply to matching entries.
type LabelRule struct {
	Field string
	Match string
	Label string
}

// Labeler attaches free-form string labels to port entries based on rules.
type Labeler struct {
	rules []LabelRule
}

// NewLabeler creates a Labeler from the provided rules.
// Returns an error if any rule has an unsupported field, empty match, or empty label.
func NewLabeler(rules []LabelRule) (*Labeler, error) {
	validFields := map[string]bool{
		"protocol": true,
		"state":    true,
		"process":  true,
		"port":     true,
	}
	for _, r := range rules {
		if !validFields[r.Field] {
			return nil, fmt.Errorf("labeler: unsupported field %q", r.Field)
		}
		if r.Match == "" {
			return nil, fmt.Errorf("labeler: match value must not be empty")
		}
		if r.Label == "" {
			return nil, fmt.Errorf("labeler: label must not be empty")
		}
	}
	return &Labeler{rules: rules}, nil
}

// Apply iterates over entries and attaches labels according to the rules.
// Multiple matching rules append multiple labels separated by commas.
func (l *Labeler) Apply(entries []PortEntry) []PortEntry {
	out := make([]PortEntry, len(entries))
	for i, e := range entries {
		labels := []string{}
		for _, r := range l.rules {
			if l.fieldValue(e, r.Field) == r.Match {
				labels = append(labels, r.Label)
			}
		}
		if len(labels) > 0 {
			if e.Tags == nil {
				e.Tags = map[string]string{}
			} else {
				copied := make(map[string]string, len(e.Tags))
				for k, v := range e.Tags {
					copied[k] = v
				}
				e.Tags = copied
			}
			for _, lbl := range labels {
				e.Tags["label:"+lbl] = lbl
			}
		}
		out[i] = e
	}
	return out
}

func (l *Labeler) fieldValue(e PortEntry, field string) string {
	switch field {
	case "protocol":
		return e.Protocol
	case "state":
		return e.State
	case "process":
		return e.Process
	case "port":
		return fmt.Sprintf("%d", e.Port)
	}
	return ""
}
