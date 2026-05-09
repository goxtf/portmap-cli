package ports

import (
	"testing"
)

var samplerEntries = []PortEntry{
	{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 443, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 8080, Protocol: "tcp", State: "LISTEN", Process: "app"},
	{Port: 5432, Protocol: "tcp", State: "LISTEN", Process: "postgres"},
	{Port: 6379, Protocol: "tcp", State: "LISTEN", Process: "redis"},
}

func TestNewSamplerInvalidN(t *testing.T) {
	_, err := NewSampler(0, "first", 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNewSamplerInvalidMode(t *testing.T) {
	_, err := NewSampler(2, "middle", 0)
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestNewSamplerValid(t *testing.T) {
	for _, mode := range []string{"random", "first", "last"} {
		_, err := NewSampler(3, mode, 42)
		if err != nil {
			t.Fatalf("unexpected error for mode %q: %v", mode, err)
		}
	}
}

func TestSampleFirst(t *testing.T) {
	s, _ := NewSampler(2, "first", 0)
	result := s.Sample(samplerEntries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Port != 80 || result[1].Port != 443 {
		t.Errorf("expected first two entries, got ports %d and %d", result[0].Port, result[1].Port)
	}
}

func TestSampleLast(t *testing.T) {
	s, _ := NewSampler(2, "last", 0)
	result := s.Sample(samplerEntries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Port != 5432 || result[1].Port != 6379 {
		t.Errorf("expected last two entries, got ports %d and %d", result[0].Port, result[1].Port)
	}
}

func TestSampleRandom(t *testing.T) {
	s, _ := NewSampler(3, "random", 99)
	result := s.Sample(samplerEntries)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestSampleNLargerThanInput(t *testing.T) {
	s, _ := NewSampler(100, "first", 0)
	result := s.Sample(samplerEntries)
	if len(result) != len(samplerEntries) {
		t.Fatalf("expected all %d entries, got %d", len(samplerEntries), len(result))
	}
}

func TestSampleEmptyInput(t *testing.T) {
	s, _ := NewSampler(3, "random", 0)
	result := s.Sample([]PortEntry{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}

func TestSampleRandomReproducible(t *testing.T) {
	s1, _ := NewSampler(3, "random", 42)
	s2, _ := NewSampler(3, "random", 42)
	r1 := s1.Sample(samplerEntries)
	r2 := s2.Sample(samplerEntries)
	for i := range r1 {
		if r1[i].Port != r2[i].Port {
			t.Errorf("expected reproducible results at index %d: %d != %d", i, r1[i].Port, r2[i].Port)
		}
	}
}
