package ports

import "strings"

// TrimOptions controls which fields are trimmed and how.
type TrimOptions struct {
	Fields    []string
	MaxLength int // 0 means no truncation
}

// Trimmer strips leading/trailing whitespace from string fields
// and optionally truncates values to a maximum length.
type Trimmer struct {
	fields    map[string]struct{}
	maxLength int
}

var validTrimFields = map[string]struct{}{
	"process": {},
	"state":   {},
	"proto":   {},
	"address": {},
}

// NewTrimmer creates a Trimmer for the given fields.
// If fields is empty all valid fields are trimmed.
// maxLength <= 0 disables truncation.
func NewTrimmer(opts TrimOptions) (*Trimmer, error) {
	fields := opts.Fields
	if len(fields) == 0 {
		for f := range validTrimFields {
			fields = append(fields, f)
		}
	}
	for _, f := range fields {
		if _, ok := validTrimFields[f]; !ok {
			return nil, fmt.Errorf("trimmer: unknown field %q", f)
		}
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	return &Trimmer{fields: set, maxLength: opts.MaxLength}, nil
}

// Apply returns a new slice of PortEntry values with the configured
// fields trimmed and optionally truncated.
func (t *Trimmer) Apply(entries []PortEntry) []PortEntry {
	out := make([]PortEntry, len(entries))
	for i, e := range entries {
		if _, ok := t.fields["process"]; ok {
			e.Process = t.trim(e.Process)
		}
		if _, ok := t.fields["state"]; ok {
			e.State = t.trim(e.State)
		}
		if _, ok := t.fields["proto"]; ok {
			e.Proto = t.trim(e.Proto)
		}
		if _, ok := t.fields["address"]; ok {
			e.Address = t.trim(e.Address)
		}
		out[i] = e
	}
	return out
}

func (t *Trimmer) trim(s string) string {
	s = strings.TrimSpace(s)
	if t.maxLength > 0 && len(s) > t.maxLength {
		return s[:t.maxLength]
	}
	return s
}
