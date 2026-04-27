package ports

import (
	"context"
	"testing"
	"time"
)

func TestNewWatcherNilScanner(t *testing.T) {
	_, err := NewWatcher(WatcherConfig{
		OnChange: func(_, _ []PortEntry) {},
	})
	if err == nil {
		t.Fatal("expected error for nil scanner")
	}
}

func TestNewWatcherNilCallback(t *testing.T) {
	s, _ := NewScanner()
	_, err := NewWatcher(WatcherConfig{Scanner: s})
	if err == nil {
		t.Fatal("expected error for nil OnChange")
	}
}

func TestNewWatcherDefaultInterval(t *testing.T) {
	s, _ := NewScanner()
	w, err := NewWatcher(WatcherConfig{
		Scanner:  s,
		OnChange: func(_, _ []PortEntry) {},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.cfg.Interval != 2*time.Second {
		t.Errorf("expected default interval 2s, got %v", w.cfg.Interval)
	}
}

func TestEntryKey(t *testing.T) {
	e := PortEntry{Protocol: "tcp", LocalAddress: "0.0.0.0:8080", State: "LISTEN"}
	expected := "tcp|0.0.0.0:8080|LISTEN"
	if got := entryKey(e); got != expected {
		t.Errorf("entryKey = %q, want %q", got, expected)
	}
}

func TestDiffEntries(t *testing.T) {
	a := map[string]PortEntry{
		"k1": {Protocol: "tcp", LocalAddress: ":80"},
		"k2": {Protocol: "tcp", LocalAddress: ":443"},
	}
	b := map[string]PortEntry{
		"k1": {Protocol: "tcp", LocalAddress: ":80"},
	}
	result := diffEntries(a, b)
	if len(result) != 1 {
		t.Fatalf("expected 1 diff entry, got %d", len(result))
	}
	if result[0].LocalAddress != ":443" {
		t.Errorf("unexpected diff entry: %+v", result[0])
	}
}

func TestWatcherCancelledContext(t *testing.T) {
	s, _ := NewScanner()
	w, err := NewWatcher(WatcherConfig{
		Scanner:  s,
		Interval: 50 * time.Millisecond,
		OnChange: func(_, _ []PortEntry) {},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err = w.Start(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("expected context error, got: %v", err)
	}
}
