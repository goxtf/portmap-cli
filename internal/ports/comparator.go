package ports

import (
	"fmt"
	"sort"
)

// CompareResult holds the diff between two slices of PortEntry.
type CompareResult struct {
	Added   []PortEntry
	Removed []PortEntry
	Common  []PortEntry
}

// Comparator compares two sets of port entries.
type Comparator struct{}

// NewComparator creates a new Comparator.
func NewComparator() *Comparator {
	return &Comparator{}
}

// Compare returns a CompareResult describing differences between baseline and current.
func (c *Comparator) Compare(baseline, current []PortEntry) CompareResult {
	baseMap := make(map[string]PortEntry, len(baseline))
	for _, e := range baseline {
		baseMap[entryKey(e)] = e
	}

	currMap := make(map[string]PortEntry, len(current))
	for _, e := range current {
		currMap[entryKey(e)] = e
	}

	var result CompareResult

	for k, e := range currMap {
		if _, exists := baseMap[k]; exists {
			result.Common = append(result.Common, e)
		} else {
			result.Added = append(result.Added, e)
		}
	}

	for k, e := range baseMap {
		if _, exists := currMap[k]; !exists {
			result.Removed = append(result.Removed, e)
		}
	}

	sort.Slice(result.Added, func(i, j int) bool { return result.Added[i].Port < result.Added[j].Port })
	sort.Slice(result.Removed, func(i, j int) bool { return result.Removed[i].Port < result.Removed[j].Port })
	sort.Slice(result.Common, func(i, j int) bool { return result.Common[i].Port < result.Common[j].Port })

	return result
}

// Summary returns a human-readable summary of the CompareResult.
func (r *CompareResult) Summary() string {
	return fmt.Sprintf("Added: %d, Removed: %d, Common: %d",
		len(r.Added), len(r.Removed), len(r.Common))
}
