package ports

import (
	"testing"
)

var splitterEntries = []PortEntry{
	{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
	{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
	{Protocol: "udp", Port: 53, State: "UNCONN", Process: "systemd-resolved"},
	{Protocol: "tcp", Port: 8080, State: "ESTABLISHED", Process: "node"},
	{Protocol: "udp", Port: 123, State: "UNCONN", Process: "ntpd"},
}

func TestNewSplitterInvalidField(t *testing.T) {
	_, err := NewSplitter("port", "80")
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestNewSplitterEmptyValue(t *testing.T) {
	_, err := NewSplitter("protocol", "")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestNewSplitterValid(t *testing.T) {
	s, err := NewSplitter("protocol", "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil splitter")
	}
}

func TestSplitByProtocol(t *testing.T) {
	s, _ := NewSplitter("protocol", "tcp")
	result := s.Split(splitterEntries)
	if len(result.Matched) != 3 {
		t.Errorf("expected 3 matched, got %d", len(result.Matched))
	}
	if len(result.Unmatched) != 2 {
		t.Errorf("expected 2 unmatched, got %d", len(result.Unmatched))
	}
}

func TestSplitByState(t *testing.T) {
	s, _ := NewSplitter("state", "LISTEN")
	result := s.Split(splitterEntries)
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(result.Matched))
	}
	if len(result.Unmatched) != 3 {
		t.Errorf("expected 3 unmatched, got %d", len(result.Unmatched))
	}
}

func TestSplitByProcess(t *testing.T) {
	s, _ := NewSplitter("process", "nginx")
	result := s.Split(splitterEntries)
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(result.Matched))
	}
}

func TestSplitEmptyInput(t *testing.T) {
	s, _ := NewSplitter("protocol", "tcp")
	result := s.Split([]PortEntry{})
	if len(result.Matched) != 0 || len(result.Unmatched) != 0 {
		t.Error("expected empty results for empty input")
	}
}

func TestSplitNoMatchProducesAllUnmatched(t *testing.T) {
	s, _ := NewSplitter("process", "unknown-proc")
	result := s.Split(splitterEntries)
	if len(result.Matched) != 0 {
		t.Errorf("expected 0 matched, got %d", len(result.Matched))
	}
	if len(result.Unmatched) != len(splitterEntries) {
		t.Errorf("expected all entries unmatched")
	}
}
