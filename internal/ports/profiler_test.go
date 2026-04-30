package ports

import (
	"testing"
)

var profilerEntries = []PortEntry{
	{Protocol: "tcp", State: "LISTEN", Process: "nginx", Port: 80},
	{Protocol: "tcp", State: "LISTEN", Process: "nginx", Port: 443},
	{Protocol: "tcp", State: "ESTABLISHED", Process: "sshd", Port: 22},
	{Protocol: "udp", State: "UNCONN", Process: "chronyd", Port: 123},
	{Protocol: "tcp", State: "LISTEN", Process: "nginx", Port: 8080},
	{Protocol: "udp", State: "UNCONN", Process: "nginx", Port: 53},
}

func TestNewProfilerInvalidN(t *testing.T) {
	_, err := NewProfiler(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNewProfilerValid(t *testing.T) {
	p, err := NewProfiler(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Profiler")
	}
}

func TestProfileTotalPorts(t *testing.T) {
	p, _ := NewProfiler(5)
	profile := p.Build(profilerEntries)
	if profile.TotalPorts != len(profilerEntries) {
		t.Errorf("expected %d, got %d", len(profilerEntries), profile.TotalPorts)
	}
}

func TestProfileByProtocol(t *testing.T) {
	p, _ := NewProfiler(5)
	profile := p.Build(profilerEntries)
	if profile.ByProtocol["tcp"] != 4 {
		t.Errorf("expected tcp=4, got %d", profile.ByProtocol["tcp"])
	}
	if profile.ByProtocol["udp"] != 2 {
		t.Errorf("expected udp=2, got %d", profile.ByProtocol["udp"])
	}
}

func TestProfileByState(t *testing.T) {
	p, _ := NewProfiler(5)
	profile := p.Build(profilerEntries)
	if profile.ByState["LISTEN"] != 3 {
		t.Errorf("expected LISTEN=3, got %d", profile.ByState["LISTEN"])
	}
}

func TestProfileTopProcessesLimit(t *testing.T) {
	p, _ := NewProfiler(2)
	profile := p.Build(profilerEntries)
	if len(profile.TopProcesses) != 2 {
		t.Errorf("expected 2 top processes, got %d", len(profile.TopProcesses))
	}
	if profile.TopProcesses[0].Process != "nginx" {
		t.Errorf("expected nginx first, got %s", profile.TopProcesses[0].Process)
	}
}

func TestProfileCapturedAtUTC(t *testing.T) {
	p, _ := NewProfiler(1)
	profile := p.Build(profilerEntries)
	if profile.CapturedAt.Location().String() != "UTC" {
		t.Errorf("expected UTC, got %s", profile.CapturedAt.Location())
	}
}

func TestProfileEmptyEntries(t *testing.T) {
	p, _ := NewProfiler(3)
	profile := p.Build([]PortEntry{})
	if profile.TotalPorts != 0 {
		t.Errorf("expected 0 total ports, got %d", profile.TotalPorts)
	}
	if len(profile.TopProcesses) != 0 {
		t.Errorf("expected no top processes, got %d", len(profile.TopProcesses))
	}
}
