package ports

import (
	"testing"
)

func enricherSampleEntries() []PortEntry {
	return []PortEntry{
		{LocalPort: 80, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
		{LocalPort: 443, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
		{LocalPort: 5432, Protocol: "TCP", Process: "postgres", State: "LISTEN"},
		{LocalPort: 9999, Protocol: "UDP", Process: "custom", State: "UNCONN"},
	}
}

func TestNewEnricherNilResolver(t *testing.T) {
	_, err := NewEnricher(nil, "")
	if err == nil {
		t.Fatal("expected error for nil resolver")
	}
}

func TestNewEnricherInvalidGroupBy(t *testing.T) {
	r, _ := NewResolver(nil)
	_, err := NewEnricher(r, "invalid_field")
	if err == nil {
		t.Fatal("expected error for invalid groupBy")
	}
}

func TestNewEnricherValid(t *testing.T) {
	r, _ := NewResolver(nil)
	e, err := NewEnricher(r, "process")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil enricher")
	}
}

func TestEnrichResolvesServiceName(t *testing.T) {
	r, _ := NewResolver(nil)
	e, _ := NewEnricher(r, "")
	entries := enricherSampleEntries()
	enriched := e.Enrich(entries)

	if len(enriched) != len(entries) {
		t.Fatalf("expected %d enriched entries, got %d", len(entries), len(enriched))
	}
	// port 80/TCP should resolve to "http"
	if enriched[0].ServiceName != "http" {
		t.Errorf("expected service 'http' for port 80, got %q", enriched[0].ServiceName)
	}
	// port 443/TCP should resolve to "https"
	if enriched[1].ServiceName != "https" {
		t.Errorf("expected service 'https' for port 443, got %q", enriched[1].ServiceName)
	}
}

func TestEnrichGroupByProcess(t *testing.T) {
	r, _ := NewResolver(nil)
	e, _ := NewEnricher(r, "process")
	enriched := e.Enrich(enricherSampleEntries())
	for _, en := range enriched {
		if en.Group != en.Process {
			t.Errorf("expected group %q, got %q", en.Process, en.Group)
		}
	}
}

func TestEnrichGroupByProtocol(t *testing.T) {
	r, _ := NewResolver(nil)
	e, _ := NewEnricher(r, "protocol")
	enriched := e.Enrich(enricherSampleEntries())
	for _, en := range enriched {
		if en.Group != en.Protocol {
			t.Errorf("expected group %q, got %q", en.Protocol, en.Group)
		}
	}
}

func TestEnrichEmptyInput(t *testing.T) {
	r, _ := NewResolver(nil)
	e, _ := NewEnricher(r, "state")
	enriched := e.Enrich([]PortEntry{})
	if len(enriched) != 0 {
		t.Errorf("expected empty result, got %d entries", len(enriched))
	}
}
