package ports

import (
	"fmt"
)

// Limiter restricts the number of entries returned from a result set.
type Limiter struct {
	max    int
	offset int
}

// NewLimiter creates a Limiter with the given max count and offset.
// max must be >= 1; offset must be >= 0.
func NewLimiter(max, offset int) (*Limiter, error) {
	if max < 1 {
		return nil, fmt.Errorf("limiter: max must be at least 1, got %d", max)
	}
	if offset < 0 {
		return nil, fmt.Errorf("limiter: offset must be non-negative, got %d", offset)
	}
	return &Limiter{max: max, offset: offset}, nil
}

// Apply returns a slice of entries starting at offset, up to max entries.
// If offset exceeds the length of entries, an empty slice is returned.
func (l *Limiter) Apply(entries []PortEntry) []PortEntry {
	if l.offset >= len(entries) {
		return []PortEntry{}
	}
	sliced := entries[l.offset:]
	if len(sliced) <= l.max {
		return sliced
	}
	return sliced[:l.max]
}

// PageInfo holds pagination metadata derived from a Limiter.
type PageInfo struct {
	Offset     int
	Max        int
	Total      int
	Returned   int
	HasMore    bool
}

// Info returns pagination metadata for the given total entry count.
func (l *Limiter) Info(total, returned int) PageInfo {
	return PageInfo{
		Offset:   l.offset,
		Max:      l.max,
		Total:    total,
		Returned: returned,
		HasMore:  l.offset+returned < total,
	}
}
