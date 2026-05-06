package ports

import "time"

// PruneOptions controls which entries are removed.
type PruneOptions struct {
	// RemoveClosed removes entries whose State is "CLOSE_WAIT" or "CLOSED".
	RemoveClosed bool
	// RemoveNoProcess removes entries that have no associated process name.
	RemoveNoProcess bool
	// OlderThan removes entries whose LastSeen is older than this duration.
	// A zero value disables time-based pruning.
	OlderThan time.Duration
	// now is injectable for testing.
	now func() time.Time
}

// Pruner removes stale or unwanted entries from a slice.
type Pruner struct {
	opts PruneOptions
}

// NewPruner returns a Pruner configured with the given options.
// At least one pruning criterion must be enabled, otherwise an error is
// returned.
func NewPruner(opts PruneOptions) (*Pruner, error) {
	if !opts.RemoveClosed && !opts.RemoveNoProcess && opts.OlderThan <= 0 {
		return nil, ErrNoPruneCriteria
	}
	if opts.now == nil {
		opts.now = time.Now
	}
	return &Pruner{opts: opts}, nil
}

// ErrNoPruneCriteria is returned when no pruning criteria are specified.
var ErrNoPruneCriteria = errorf("pruner: at least one pruning criterion must be enabled")

// Apply returns a new slice with matching entries removed.
func (p *Pruner) Apply(entries []PortEntry) []PortEntry {
	result := make([]PortEntry, 0, len(entries))
	now := p.opts.now()
	for _, e := range entries {
		if p.opts.RemoveClosed && isClosed(e.State) {
			continue
		}
		if p.opts.RemoveNoProcess && e.Process == "" {
			continue
		}
		if p.opts.OlderThan > 0 && !e.LastSeen.IsZero() && now.Sub(e.LastSeen) > p.opts.OlderThan {
			continue
		}
		result = append(result, e)
	}
	return result
}

func isClosed(state string) bool {
	return state == "CLOSED" || state == "CLOSE_WAIT"
}
