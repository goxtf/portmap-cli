package ports

import (
	"testing"
)

var sorterSampleEntries = []PortEntry{
	{Port: 8080, Protocol: "TCP", Process: "node", State: "LISTEN"},
	{Port: 22, Protocol: "TCP", Process: "sshd", State: "LISTEN"},
	{Port: 443, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
	{Port: 5432, Protocol: "TCP", Process: "postgres", State: "ESTABLISHED"},
	{Port: 80, Protocol: "UDP", Process: "apache", State: "LISTEN"},
}

func TestNewSorterValidFields(t *testing.T) {
	for _, field := range []string{"port", "protocol", "process", "state", "PORT", "Protocol"} {
		_, err := NewSorter(field, false)
		if err != nil {
			t.Errorf("expected no error for field %q, got %v", field, err)
		}
	}
}

func TestNewSorterInvalidField(t *testing.T) {
	_, err := NewSorter("invalid", false)
	if err == nil {
		t.Error("expected error for invalid field, got nil")
	}
}

func TestSortByPort(t *testing.T) {
	s, _ := NewSorter("port", false)
	result := s.Sort(sorterSampleEntries)
	for i := 1; i < len(result); i++ {
		if result[i].Port < result[i-1].Port {
			t.Errorf("expected ascending port order, got %d before %d", result[i-1].Port, result[i].Port)
		}
	}
}

func TestSortByPortReverse(t *testing.T) {
	s, _ := NewSorter("port", true)
	result := s.Sort(sorterSampleEntries)
	for i := 1; i < len(result); i++ {
		if result[i].Port > result[i-1].Port {
			t.Errorf("expected descending port order, got %d before %d", result[i-1].Port, result[i].Port)
		}
	}
}

func TestSortByProcess(t *testing.T) {
	s, _ := NewSorter("process", false)
	result := s.Sort(sorterSampleEntries)
	for i := 1; i < len(result); i++ {
		if result[i].Process < result[i-1].Process {
			t.Errorf("expected ascending process order, got %q before %q", result[i-1].Process, result[i].Process)
		}
	}
}

func TestSortDoesNotMutateOriginal(t *testing.T) {
	original := make([]PortEntry, len(sorterSampleEntries))
	copy(original, sorterSampleEntries)
	s, _ := NewSorter("port", false)
	s.Sort(sorterSampleEntries)
	for i, e := range sorterSampleEntries {
		if e != original[i] {
			t.Error("Sort mutated the original slice")
		}
	}
}
