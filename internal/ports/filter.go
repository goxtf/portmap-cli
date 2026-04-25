package ports

import "strings"

// FilterOptions defines criteria for filtering port mappings.
type FilterOptions struct {
	Protocol string
	State    string
	PIDStr   string
	Process  string
	Port     string
}

// Filter holds the parsed filter options and applies them to port entries.
type Filter struct {
	opts FilterOptions
}

// NewFilter creates a new Filter from the given options.
func NewFilter(opts FilterOptions) *Filter {
	return &Filter{opts: opts}
}

// Apply returns only the entries that match all non-empty filter criteria.
func (f *Filter) Apply(entries []PortEntry) []PortEntry {
	var result []PortEntry
	for _, e := range entries {
		if !f.matches(e) {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (f *Filter) matches(e PortEntry) bool {
	if f.opts.Protocol != "" && !strings.EqualFold(e.Protocol, f.opts.Protocol) {
		return false
	}
	if f.opts.State != "" && !strings.EqualFold(e.State, f.opts.State) {
		return false
	}
	if f.opts.PIDStr != "" && e.PID != f.opts.PIDStr {
		return false
	}
	if f.opts.Process != "" && !strings.Contains(strings.ToLower(e.Process), strings.ToLower(f.opts.Process)) {
		return false
	}
	if f.opts.Port != "" && e.LocalPort != f.opts.Port {
		return false
	}
	return true
}
