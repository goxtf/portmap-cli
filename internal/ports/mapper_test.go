package ports

import (
	"testing"
)

var mapperEntries = []PortEntry{
	{Port: 80, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
	{Port: 80, Protocol: "TCP", Process: "apache", State: "LISTEN"},
	{Port: 443, Protocol: "TCP", Process: "nginx", State: "LISTEN"},
	{Port: 8080, Protocol: "TCP", Process: "app", State: "CLOSE_WAIT"},
	{Port: 22, Protocol: "TCP", Process: "sshd", State: "LISTEN"},
}

func TestNewMapper(t *testing.T) {
	m, err := NewMapper(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Mapper")
	}
}

func TestMapperBuildExcludesClosed(t *testing.T) {
	m, _ := NewMapper(false)
	pm := m.Build(mapperEntries)
	if _, ok := pm[8080]; ok {
		t.Error("expected CLOSE_WAIT entry to be excluded")
	}
}

func TestMapperBuildIncludesClosed(t *testing.T) {
	m, _ := NewMapper(true)
	pm := m.Build(mapperEntries)
	if _, ok := pm[8080]; !ok {
		t.Error("expected CLOSE_WAIT entry to be included when includeClosed=true")
	}
}

func TestMapperPorts(t *testing.T) {
	m, _ := NewMapper(false)
	pm := m.Build(mapperEntries)
	ports := pm.Ports()
	if len(ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(ports))
	}
	if ports[0] != 22 || ports[1] != 80 || ports[2] != 443 {
		t.Errorf("unexpected port order: %v", ports)
	}
}

func TestMapperLookup(t *testing.T) {
	m, _ := NewMapper(false)
	pm := m.Build(mapperEntries)
	entries, ok := pm.Lookup(80)
	if !ok {
		t.Fatal("expected port 80 to exist")
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries for port 80, got %d", len(entries))
	}
}

func TestMapperLookupMissing(t *testing.T) {
	m, _ := NewMapper(false)
	pm := m.Build(mapperEntries)
	_, ok := pm.Lookup(9999)
	if ok {
		t.Error("expected port 9999 to be absent")
	}
}

func TestMapperConflicts(t *testing.T) {
	m, _ := NewMapper(false)
	pm := m.Build(mapperEntries)
	conflicts := pm.Conflicts()
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}
	if _, ok := conflicts[80]; !ok {
		t.Error("expected port 80 to be in conflicts")
	}
}
