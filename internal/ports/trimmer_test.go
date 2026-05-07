package ports

import (
	"testing"
)

func trimmerEntries() []PortEntry {
	return []PortEntry{
		{Port: 80, Proto: "  tcp  ", State: " LISTEN ", Process: "  nginx  ", Address: " 0.0.0.0 "},
		{Port: 443, Proto: "tcp", State: "ESTABLISHED", Process: "caddy", Address: "127.0.0.1"},
		{Port: 8080, Proto: " udp ", State: " CLOSED ", Process: " long-process-name ", Address: " ::1 "},
	}
}

func TestNewTrimmerAllFields(t *testing.T) {
	_, err := NewTrimmer(TrimOptions{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewTrimmerInvalidField(t *testing.T) {
	_, err := NewTrimmer(TrimOptions{Fields: []string{"invalid"}})
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestNewTrimmerValidFields(t *testing.T) {
	_, err := NewTrimmer(TrimOptions{Fields: []string{"process", "state"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTrimStripsWhitespace(t *testing.T) {
	tr, _ := NewTrimmer(TrimOptions{})
	out := tr.Apply(trimmerEntries())
	if out[0].Proto != "tcp" {
		t.Errorf("expected 'tcp', got %q", out[0].Proto)
	}
	if out[0].State != "LISTEN" {
		t.Errorf("expected 'LISTEN', got %q", out[0].State)
	}
	if out[0].Process != "nginx" {
		t.Errorf("expected 'nginx', got %q", out[0].Process)
	}
	if out[0].Address != "0.0.0.0" {
		t.Errorf("expected '0.0.0.0', got %q", out[0].Address)
	}
}

func TestTrimMaxLength(t *testing.T) {
	tr, _ := NewTrimmer(TrimOptions{Fields: []string{"process"}, MaxLength: 4})
	out := tr.Apply(trimmerEntries())
	if out[2].Process != "long" {
		t.Errorf("expected 'long', got %q", out[2].Process)
	}
	// other fields should be unchanged (still have spaces)
	if out[2].Proto != " udp " {
		t.Errorf("expected ' udp ', got %q", out[2].Proto)
	}
}

func TestTrimDoesNotMutateOriginal(t *testing.T) {
	entries := trimmerEntries()
	origProto := entries[0].Proto
	tr, _ := NewTrimmer(TrimOptions{})
	tr.Apply(entries)
	if entries[0].Proto != origProto {
		t.Error("original entries were mutated")
	}
}

func TestTrimAlreadyClean(t *testing.T) {
	tr, _ := NewTrimmer(TrimOptions{})
	out := tr.Apply(trimmerEntries())
	if out[1].Proto != "tcp" {
		t.Errorf("expected 'tcp', got %q", out[1].Proto)
	}
	if out[1].Process != "caddy" {
		t.Errorf("expected 'caddy', got %q", out[1].Process)
	}
}
