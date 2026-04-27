package ports

import (
	"context"
	"time"
)

// WatcherConfig holds configuration for the port watcher.
type WatcherConfig struct {
	Interval  time.Duration
	Scanner   *Scanner
	OnChange  func(added, removed []PortEntry)
}

// Watcher periodically scans ports and emits diffs.
type Watcher struct {
	cfg      WatcherConfig
	previous map[string]PortEntry
}

// NewWatcher creates a new Watcher with the given config.
// Returns an error if the scanner or onChange callback is nil.
func NewWatcher(cfg WatcherConfig) (*Watcher, error) {
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("watcher: scanner must not be nil")
	}
	if cfg.OnChange == nil {
		return nil, fmt.Errorf("watcher: OnChange callback must not be nil")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 2 * time.Second
	}
	return &Watcher{
		cfg:      cfg,
		previous: make(map[string]PortEntry),
	}, nil
}

// Start begins watching for port changes until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	// Perform an initial scan to populate baseline.
	if err := w.tick(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(); err != nil {
				return err
			}
		}
	}
}

func (w *Watcher) tick() error {
	entries, err := w.cfg.Scanner.Scan()
	if err != nil {
		return err
	}

	current := make(map[string]PortEntry, len(entries))
	for _, e := range entries {
		current[entryKey(e)] = e
	}

	added := diffEntries(current, w.previous)
	removed := diffEntries(w.previous, current)

	if len(added) > 0 || len(removed) > 0 {
		w.cfg.OnChange(added, removed)
	}

	w.previous = current
	return nil
}

// entryKey returns a unique string key for a PortEntry.
func entryKey(e PortEntry) string {
	return e.Protocol + "|" + e.LocalAddress + "|" + e.State
}

// diffEntries returns entries present in a but not in b.
func diffEntries(a, b map[string]PortEntry) []PortEntry {
	var result []PortEntry
	for k, v := range a {
		if _, ok := b[k]; !ok {
			result = append(result, v)
		}
	}
	return result
}
