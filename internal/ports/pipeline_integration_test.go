package ports_test

import (
	"testing"

	"github.com/user/portmap-cli/internal/ports"
)

func integrationPipelineEntries() []ports.PortEntry {
	return []ports.PortEntry{
		{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
		{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
		{Protocol: "udp", Port: 53, State: "UNCONN", Process: "dnsmasq"},
		{Protocol: "tcp", Port: 8080, State: "LISTEN", Process: "app"},
		{Protocol: "tcp", Port: 22, State: "LISTEN", Process: "sshd"},
		{Protocol: "tcp", Port: 22, State: "LISTEN", Process: "sshd"},
	}
}

func TestPipelineFilterThenDeduplicate(t *testing.T) {
	f, _ := ports.NewFilter("tcp", "", "")
	d, _ := ports.NewDeduplicator([]string{"port", "protocol", "process"})

	p, err := ports.NewPipeline(
		func(e []ports.PortEntry) ([]ports.PortEntry, error) { return f.Apply(e), nil },
		func(e []ports.PortEntry) ([]ports.PortEntry, error) { return d.Deduplicate(e), nil },
	)
	if err != nil {
		t.Fatalf("unexpected error building pipeline: %v", err)
	}

	result, err := p.Run(integrationPipelineEntries())
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	// 4 unique tcp entries (80 nginx, 443 nginx, 8080 app, 22 sshd)
	if len(result) != 4 {
		t.Fatalf("expected 4 entries after filter+dedup, got %d", len(result))
	}
	for _, e := range result {
		if e.Protocol != "tcp" {
			t.Errorf("non-tcp entry slipped through: %+v", e)
		}
	}
}

func TestPipelineFilterSortExport(t *testing.T) {
	f, _ := ports.NewFilter("tcp", "", "")
	s, _ := ports.NewSorter("port", false)

	p, err := ports.NewPipeline(
		func(e []ports.PortEntry) ([]ports.PortEntry, error) { return f.Apply(e), nil },
		func(e []ports.PortEntry) ([]ports.PortEntry, error) { return s.Sort(e), nil },
	)
	if err != nil {
		t.Fatalf("pipeline build error: %v", err)
	}

	result, err := p.Run(integrationPipelineEntries())
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	for i := 1; i < len(result); i++ {
		if result[i].Port < result[i-1].Port {
			t.Errorf("entries not sorted at index %d: %d < %d", i, result[i].Port, result[i-1].Port)
		}
	}
}
