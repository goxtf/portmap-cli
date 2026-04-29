package ports

import (
	"fmt"
	"strings"
)

// EnrichedEntry wraps a PortEntry with additional resolved metadata.
type EnrichedEntry struct {
	PortEntry
	ServiceName string `json:"service_name"`
	Group       string `json:"group"`
}

// Enricher combines resolver and grouper to produce enriched port entries.
type Enricher struct {
	resolver *Resolver
	groupBy  string
}

// NewEnricher creates an Enricher using the provided resolver and group field.
// Valid groupBy values: "process", "protocol", "state", "" (none).
func NewEnricher(resolver *Resolver, groupBy string) (*Enricher, error) {
	if resolver == nil {
		return nil, fmt.Errorf("enricher: resolver must not be nil")
	}
	valid := map[string]bool{"process": true, "protocol": true, "state": true, "": true}
	if !valid[strings.ToLower(groupBy)] {
		return nil, fmt.Errorf("enricher: invalid groupBy field %q", groupBy)
	}
	return &Enricher{resolver: resolver, groupBy: strings.ToLower(groupBy)}, nil
}

// Enrich resolves service names and assigns group labels to each entry.
func (e *Enricher) Enrich(entries []PortEntry) []EnrichedEntry {
	result := make([]EnrichedEntry, 0, len(entries))
	for _, entry := range entries {
		svc := e.resolver.Resolve(entry.LocalPort, entry.Protocol)
		group := e.computeGroup(entry)
		result = append(result, EnrichedEntry{
			PortEntry:   entry,
			ServiceName: svc,
			Group:       group,
		})
	}
	return result
}

func (e *Enricher) computeGroup(entry PortEntry) string {
	switch e.groupBy {
	case "process":
		return entry.Process
	case "protocol":
		return entry.Protocol
	case "state":
		return entry.State
	default:
		return ""
	}
}
