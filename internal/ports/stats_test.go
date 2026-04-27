package ports

import (
	"testing"
)

func statsEntries() []PortEntry {
	return []PortEntry{
		{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
		{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
		{Protocol: "tcp", Port: 8080, State: "ESTABLISHED", Process: "go"},
		{Protocol: "udp", Port: 53, State: "UNCONN", Process: "systemd-resolved"},
		{Protocol: "tcp", Port: 22, State: "LISTEN", Process: "sshd"},
		{Protocol: "tcp", Port: 5432, State: "LISTEN", Process: "postgres"},
	}
}

func TestNewStatsCollector(t *testing.T) {
	sc := NewStatsCollector()
	if sc == nil {
		t.Fatal("expected non-nil StatsCollector")
	}
}

func TestCollectTotalEntries(t *testing.T) {
	sc := NewStatsCollector()
	stats := sc.Collect(statsEntries())
	if stats.TotalEntries != 6 {
		t.Errorf("expected 6, got %d", stats.TotalEntries)
	}
}

func TestCollectByProtocol(t *testing.T) {
	sc := NewStatsCollector()
	stats := sc.Collect(statsEntries())
	if stats.ByProtocol["tcp"] != 5 {
		t.Errorf("expected tcp=5, got %d", stats.ByProtocol["tcp"])
	}
	if stats.ByProtocol["udp"] != 1 {
		t.Errorf("expected udp=1, got %d", stats.ByProtocol["udp"])
	}
}

func TestCollectByState(t *testing.T) {
	sc := NewStatsCollector()
	stats := sc.Collect(statsEntries())
	if stats.ByState["LISTEN"] != 4 {
		t.Errorf("expected LISTEN=4, got %d", stats.ByState["LISTEN"])
	}
}

func TestCollectUniqueProcesses(t *testing.T) {
	sc := NewStatsCollector()
	stats := sc.Collect(statsEntries())
	if stats.UniqueProcesses != 5 {
		t.Errorf("expected 5 unique processes, got %d", stats.UniqueProcesses)
	}
}

func TestCollectUniquePorts(t *testing.T) {
	sc := NewStatsCollector()
	stats := sc.Collect(statsEntries())
	if stats.UniquePorts != 6 {
		t.Errorf("expected 6 unique ports, got %d", stats.UniquePorts)
	}
}

func TestCollectTopProcesses(t *testing.T) {
	sc := NewStatsCollector()
	stats := sc.Collect(statsEntries())
	if len(stats.TopProcesses) == 0 {
		t.Fatal("expected non-empty TopProcesses")
	}
	if stats.TopProcesses[0].Process != "nginx" {
		t.Errorf("expected nginx at top, got %s", stats.TopProcesses[0].Process)
	}
	if stats.TopProcesses[0].Count != 2 {
		t.Errorf("expected count 2, got %d", stats.TopProcesses[0].Count)
	}
}

func TestCollectEmptyEntries(t *testing.T) {
	sc := NewStatsCollector()
	stats := sc.Collect([]PortEntry{})
	if stats.TotalEntries != 0 {
		t.Errorf("expected 0, got %d", stats.TotalEntries)
	}
	if len(stats.TopProcesses) != 0 {
		t.Errorf("expected empty TopProcesses")
	}
}
