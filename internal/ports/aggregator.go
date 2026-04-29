package ports

import "sort"

// AggregateResult holds aggregated statistics for a group of entries.
type AggregateResult struct {
	Key      string
	Count    int
	Ports    []int
	Protocols []string
}

// Aggregator combines grouping and counting of port entries.
type Aggregator struct {
	field string
}

// NewAggregator returns an Aggregator for the given field.
// Valid fields: "process", "protocol", "state".
func NewAggregator(field string) (*Aggregator, error) {
	valid := map[string]bool{"process": true, "protocol": true, "state": true}
	if !valid[field] {
		return nil, &InvalidFieldError{Field: field}
	}
	return &Aggregator{field: field}, nil
}

// Aggregate groups entries by the configured field and returns sorted results.
func (a *Aggregator) Aggregate(entries []PortEntry) []AggregateResult {
	groups := make(map[string][]PortEntry)
	for _, e := range entries {
		key := a.keyFor(e)
		groups[key] = append(groups[key], e)
	}

	results := make([]AggregateResult, 0, len(groups))
	for key, group := range groups {
		results = append(results, a.buildResult(key, group))
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Key < results[j].Key
	})
	return results
}

func (a *Aggregator) keyFor(e PortEntry) string {
	switch a.field {
	case "process":
		return e.Process
	case "protocol":
		return e.Protocol
	case "state":
		return e.State
	}
	return ""
}

func (a *Aggregator) buildResult(key string, entries []PortEntry) AggregateResult {
	portSet := make(map[int]bool)
	protoSet := make(map[string]bool)
	for _, e := range entries {
		portSet[e.Port] = true
		protoSet[e.Protocol] = true
	}
	ports := make([]int, 0, len(portSet))
	for p := range portSet {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	protos := make([]string, 0, len(protoSet))
	for p := range protoSet {
		protos = append(protos, p)
	}
	sort.Strings(protos)
	return AggregateResult{
		Key:       key,
		Count:     len(entries),
		Ports:     ports,
		Protocols: protos,
	}
}

// InvalidFieldError is returned when an unsupported aggregation field is used.
type InvalidFieldError struct {
	Field string
}

func (e *InvalidFieldError) Error() string {
	return "aggregator: invalid field: " + e.Field
}
