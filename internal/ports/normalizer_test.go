package ports

import (
	"testing"
)

var normalizerSampleEntries = []PortEntry{
	{Protocol: "TCP", State: "listen", Process: "  nginx  ", Port: 80},
	{Protocol: "UDP", State: "established", Process: "go", Port: 9000},
	{Protocol: "tcp", State: "LISTEN", Process: "python3", Port: 8080},
}

func TestNewNormalizerAllFields(t *testing.T) {
	n, err := NewNormalizer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil normalizer")
	}
}

func TestNewNormalizerInvalidField(t *testing.T) {
	_, err := NewNormalizer([]NormalizeField{"invalid"})
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestNewNormalizerValidFields(t *testing.T) {
	n, err := NewNormalizer([]NormalizeField{NormalizeProtocol, NormalizeState})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil normalizer")
	}
}

func TestNormalizeProtocol(t *testing.T) {
	n, _ := NewNormalizer([]NormalizeField{NormalizeProtocol})
	out := n.Normalize(normalizerSampleEntries)
	for _, e := range out {
		if e.Protocol != "tcp" && e.Protocol != "udp" {
			t.Errorf("expected lowercase protocol, got %q", e.Protocol)
		}
	}
}

func TestNormalizeState(t *testing.T) {
	n, _ := NewNormalizer([]NormalizeField{NormalizeState})
	out := n.Normalize(normalizerSampleEntries)
	for _, e := range out {
		if e.State != "LISTEN" && e.State != "ESTABLISHED" {
			t.Errorf("expected uppercase state, got %q", e.State)
		}
	}
}

func TestNormalizeProcessTrimsSpace(t *testing.T) {
	n, _ := NewNormalizer([]NormalizeField{NormalizeProcess})
	out := n.Normalize(normalizerSampleEntries)
	if out[0].Process != "nginx" {
		t.Errorf("expected trimmed process, got %q", out[0].Process)
	}
}

func TestNormalizeDoesNotMutateOriginal(t *testing.T) {
	original := []PortEntry{
		{Protocol: "TCP", State: "listen", Process: "  sshd  ", Port: 22},
	}
	n, _ := NewNormalizer(nil)
	n.Normalize(original)
	if original[0].Protocol != "TCP" {
		t.Error("original entry was mutated")
	}
}
