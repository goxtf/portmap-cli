package ports

import (
	"fmt"
	"sort"
)

// TraceResult holds the trace output for a single port.
type TraceResult struct {
	Port     int
	Protocol string
	Process  string
	State    string
	Hops     []string
	Notes    []string
}

// Tracer follows a port through filter, tag, and label metadata to
// produce a human-readable audit trail of how an entry was processed.
type Tracer struct {
	resolver *Resolver
}

// NewTracer creates a Tracer. resolver may be nil, in which case
// service name resolution is skipped.
func NewTracer(resolver *Resolver) *Tracer {
	return &Tracer{resolver: resolver}
}

// Trace inspects every entry whose port matches the given port number
// and builds a TraceResult describing it.
func (t *Tracer) Trace(entries []PortEntry, port int) ([]TraceResult, error) {
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("tracer: invalid port %d", port)
	}

	var results []TraceResult
	for _, e := range entries {
		if e.Port != port {
			continue
		}
		r := TraceResult{
			Port:     e.Port,
			Protocol: e.Protocol,
			Process:  e.Process,
			State:    e.State,
		}

		r.Hops = append(r.Hops, fmt.Sprintf("port=%d matched", port))
		r.Hops = append(r.Hops, fmt.Sprintf("protocol=%s", e.Protocol))
		r.Hops = append(r.Hops, fmt.Sprintf("state=%s", e.State))

		if e.Process != "" {
			r.Hops = append(r.Hops, fmt.Sprintf("process=%s", e.Process))
		} else {
			r.Notes = append(r.Notes, "no owning process detected")
		}

		if t.resolver != nil {
			svc := t.resolver.Resolve(e.Port)
			if svc != "" {
				r.Hops = append(r.Hops, fmt.Sprintf("service=%s", svc))
			}
		}

		for k, v := range e.Tags {
			keys := make([]string, 0)
			_ = k
			_ = v
			_ = keys
			break
		}
		tagKeys := make([]string, 0, len(e.Tags))
		for k := range e.Tags {
			tagKeys = append(tagKeys, k)
		}
		sort.Strings(tagKeys)
		for _, k := range tagKeys {
			r.Hops = append(r.Hops, fmt.Sprintf("tag:%s=%s", k, e.Tags[k]))
		}

		if e.State == "CLOSE_WAIT" || e.State == "CLOSED" {
			r.Notes = append(r.Notes, "connection is closing or already closed")
		}

		results = append(results, r)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("tracer: no entries found for port %d", port)
	}
	return results, nil
}
