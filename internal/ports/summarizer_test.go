package ports

import (
	"testing"
)

var summarizerEntries = []PortEntry{
	{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
	{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
	{Protocol: "tcp", Port: 8080, State: "ESTABLISHED", Process: "go"},
	{Protocol: "udp", Port: 53, State: "LISTEN", Process: "dnsmasq"},
	{Protocol: "tcp", Port: 22, State: "LISTEN", Process: "sshd"},
	{Protocol: "tcp", Port: 80, State: "ESTABLISHED", Process: "nginx"},
}

func TestNewSummarizerInvalidN(t *testing.T) {
	_, err := NewSummarizer(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNewSummarizerValid(t *testing.T) {
	s, err := NewSummarizer(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil summarizer")
	}
}

func TestSummarizeTotal(t *testing.T) {
	s, _ := NewSummarizer(5)
	sum := s.Summarize(summarizerEntries)
	if sum.Total != len(summarizerEntries) {
		t.Errorf("expected total %d, got %d", len(summarizerEntries), sum.Total)
	}
}

func TestSummarizeByProtocol(t *testing.T) {
	s, _ := NewSummarizer(5)
	sum := s.Summarize(summarizerEntries)
	if sum.ByProtocol["tcp"] != 5 {
		t.Errorf("expected 5 tcp entries, got %d", sum.ByProtocol["tcp"])
	}
	if sum.ByProtocol["udp"] != 1 {
		t.Errorf("expected 1 udp entry, got %d", sum.ByProtocol["udp"])
	}
}

func TestSummarizeByState(t *testing.T) {
	s, _ := NewSummarizer(5)
	sum := s.Summarize(summarizerEntries)
	if sum.ByState["LISTEN"] != 4 {
		t.Errorf("expected 4 LISTEN, got %d", sum.ByState["LISTEN"])
	}
}

func TestSummarizeUniquePorts(t *testing.T) {
	s, _ := NewSummarizer(5)
	sum := s.Summarize(summarizerEntries)
	// ports: 80, 443, 8080, 53, 22 => 5 unique
	if sum.UniquePorts != 5 {
		t.Errorf("expected 5 unique ports, got %d", sum.UniquePorts)
	}
}

func TestSummarizeTopProcesses(t *testing.T) {
	s, _ := NewSummarizer(1)
	sum := s.Summarize(summarizerEntries)
	if len(sum.TopProcesses) != 1 {
		t.Fatalf("expected 1 top process, got %d", len(sum.TopProcesses))
	}
	if sum.TopProcesses[0] != "nginx" {
		t.Errorf("expected top process nginx, got %s", sum.TopProcesses[0])
	}
}

func TestSummarizeEmptyEntries(t *testing.T) {
	s, _ := NewSummarizer(3)
	sum := s.Summarize([]PortEntry{})
	if sum.Total != 0 {
		t.Errorf("expected total 0, got %d", sum.Total)
	}
	if len(sum.TopProcesses) != 0 {
		t.Errorf("expected no top processes, got %v", sum.TopProcesses)
	}
}
