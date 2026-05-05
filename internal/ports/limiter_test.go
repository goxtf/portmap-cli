package ports

import (
	"testing"
)

func makeLimiterEntries(n int) []PortEntry {
	entries := make([]PortEntry, n)
	for i := range entries {
		entries[i] = PortEntry{
			Port:     uint16(8000 + i),
			Protocol: "tcp",
			Process:  "proc",
			State:    "LISTEN",
		}
	}
	return entries
}

func TestNewLimiterInvalidMax(t *testing.T) {
	_, err := NewLimiter(0, 0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewLimiterNegativeOffset(t *testing.T) {
	_, err := NewLimiter(10, -1)
	if err == nil {
		t.Fatal("expected error for negative offset")
	}
}

func TestNewLimiterValid(t *testing.T) {
	l, err := NewLimiter(5, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.max != 5 || l.offset != 2 {
		t.Errorf("unexpected limiter state: max=%d offset=%d", l.max, l.offset)
	}
}

func TestLimiterApplyBasic(t *testing.T) {
	l, _ := NewLimiter(3, 0)
	entries := makeLimiterEntries(10)
	result := l.Apply(entries)
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestLimiterApplyWithOffset(t *testing.T) {
	l, _ := NewLimiter(3, 7)
	entries := makeLimiterEntries(10)
	result := l.Apply(entries)
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
	if result[0].Port != 8007 {
		t.Errorf("expected first port 8007, got %d", result[0].Port)
	}
}

func TestLimiterApplyOffsetBeyondLength(t *testing.T) {
	l, _ := NewLimiter(5, 20)
	entries := makeLimiterEntries(10)
	result := l.Apply(entries)
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestLimiterApplyFewerThanMax(t *testing.T) {
	l, _ := NewLimiter(100, 0)
	entries := makeLimiterEntries(5)
	result := l.Apply(entries)
	if len(result) != 5 {
		t.Errorf("expected 5 entries, got %d", len(result))
	}
}

func TestLimiterInfoHasMore(t *testing.T) {
	l, _ := NewLimiter(5, 0)
	info := l.Info(20, 5)
	if !info.HasMore {
		t.Error("expected HasMore=true")
	}
	if info.Total != 20 || info.Returned != 5 {
		t.Errorf("unexpected info: %+v", info)
	}
}

func TestLimiterInfoNoMore(t *testing.T) {
	l, _ := NewLimiter(5, 15)
	info := l.Info(20, 5)
	if info.HasMore {
		t.Error("expected HasMore=false")
	}
}
