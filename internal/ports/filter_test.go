package ports

import (
	"testing"
)

func sampleEntries() []PortEntry {
	return []PortEntry{
		{Protocol: "TCP", LocalPort: "8080", State: "LISTEN", PID: "123", Process: "nginx"},
		{Protocol: "UDP", LocalPort: "53", State: "", PID: "456", Process: "dnsmasq"},
		{Protocol: "TCP", LocalPort: "443", State: "LISTEN", PID: "789", Process: "nginx"},
		{Protocol: "TCP", LocalPort: "5432", State: "ESTABLISHED", PID: "321", Process: "postgres"},
	}
}

func TestFilterByProtocol(t *testing.T) {
	f := NewFilter(FilterOptions{Protocol: "UDP"})
	result := f.Apply(sampleEntries())
	if len(result) != 1 || result[0].LocalPort != "53" {
		t.Errorf("expected 1 UDP entry on port 53, got %+v", result)
	}
}

func TestFilterByState(t *testing.T) {
	f := NewFilter(FilterOptions{State: "LISTEN"})
	result := f.Apply(sampleEntries())
	if len(result) != 2 {
		t.Errorf("expected 2 LISTEN entries, got %d", len(result))
	}
}

func TestFilterByProcess(t *testing.T) {
	f := NewFilter(FilterOptions{Process: "nginx"})
	result := f.Apply(sampleEntries())
	if len(result) != 2 {
		t.Errorf("expected 2 nginx entries, got %d", len(result))
	}
}

func TestFilterByPort(t *testing.T) {
	f := NewFilter(FilterOptions{Port: "5432"})
	result := f.Apply(sampleEntries())
	if len(result) != 1 || result[0].Process != "postgres" {
		t.Errorf("expected 1 postgres entry, got %+v", result)
	}
}

func TestFilterByPID(t *testing.T) {
	f := NewFilter(FilterOptions{PIDStr: "456"})
	result := f.Apply(sampleEntries())
	if len(result) != 1 || result[0].Process != "dnsmasq" {
		t.Errorf("expected 1 dnsmasq entry, got %+v", result)
	}
}

func TestFilterNoMatch(t *testing.T) {
	f := NewFilter(FilterOptions{Protocol: "TCP", State: "LISTEN", Process: "postgres"})
	result := f.Apply(sampleEntries())
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestFilterEmpty(t *testing.T) {
	f := NewFilter(FilterOptions{})
	result := f.Apply(sampleEntries())
	if len(result) != len(sampleEntries()) {
		t.Errorf("expected all entries with empty filter, got %d", len(result))
	}
}
