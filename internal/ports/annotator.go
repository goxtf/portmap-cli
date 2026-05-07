package ports

import "fmt"

// AnnotationRule maps a field/value pair to a free-form note.
type AnnotationRule struct {
	Field string
	Match string
	Note  string
}

// Annotator attaches human-readable notes to port entries based on rules.
type Annotator struct {
	rules []AnnotationRule
}

var validAnnotatorFields = map[string]struct{}{
	"protocol": {},
	"state":    {},
	"process":  {},
	"port":     {},
}

// NewAnnotator validates rules and returns an Annotator.
func NewAnnotator(rules []AnnotationRule) (*Annotator, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("annotator: at least one rule is required")
	}
	for _, r := range rules {
		if _, ok := validAnnotatorFields[r.Field]; !ok {
			return nil, fmt.Errorf("annotator: invalid field %q", r.Field)
		}
		if r.Match == "" {
			return nil, fmt.Errorf("annotator: match value must not be empty")
		}
		if r.Note == "" {
			return nil, fmt.Errorf("annotator: note must not be empty")
		}
	}
	return &Annotator{rules: rules}, nil
}

// Annotate returns a copy of entries with Notes populated according to rules.
func (a *Annotator) Annotate(entries []PortEntry) []PortEntry {
	out := make([]PortEntry, len(entries))
	for i, e := range entries {
		copy := e
		for _, r := range a.rules {
			if a.fieldValue(copy, r.Field) == r.Match {
				if copy.Notes == "" {
					copy.Notes = r.Note
				} else {
					copy.Notes += "; " + r.Note
				}
			}
		}
		out[i] = copy
	}
	return out
}

func (a *Annotator) fieldValue(e PortEntry, field string) string {
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
