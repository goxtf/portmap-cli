package ports

import (
	"errors"
	"testing"
)

func pipelineEntries() []PortEntry {
	return []PortEntry{
		{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
		{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
		{Protocol: "udp", Port: 53, State: "UNCONN", Process: "dnsmasq"},
		{Protocol: "tcp", Port: 22, State: "LISTEN", Process: "sshd"},
	}
}

func TestNewPipelineNoSteps(t *testing.T) {
	_, err := NewPipeline()
	if err == nil {
		t.Fatal("expected error for empty steps")
	}
}

func TestNewPipelineValid(t *testing.T) {
	p, err := NewPipeline(func(e []PortEntry) ([]PortEntry, error) { return e, nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestPipelineRunSingleStep(t *testing.T) {
	filterTCP := func(entries []PortEntry) ([]PortEntry, error) {
		var out []PortEntry
		for _, e := range entries {
			if e.Protocol == "tcp" {
				out = append(out, e)
			}
		}
		return out, nil
	}
	p, _ := NewPipeline(filterTCP)
	result, err := p.Run(pipelineEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 tcp entries, got %d", len(result))
	}
}

func TestPipelineRunMultipleSteps(t *testing.T) {
	filterTCP := func(entries []PortEntry) ([]PortEntry, error) {
		var out []PortEntry
		for _, e := range entries {
			if e.Protocol == "tcp" {
				out = append(out, e)
			}
		}
		return out, nil
	}
	filterNginx := func(entries []PortEntry) ([]PortEntry, error) {
		var out []PortEntry
		for _, e := range entries {
			if e.Process == "nginx" {
				out = append(out, e)
			}
		}
		return out, nil
	}
	p, _ := NewPipeline(filterTCP, filterNginx)
	result, err := p.Run(pipelineEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 nginx/tcp entries, got %d", len(result))
	}
}

func TestPipelineRunStopOnError(t *testing.T) {
	sentinel := errors.New("step failed")
	called := false
	failStep := func(entries []PortEntry) ([]PortEntry, error) {
		return nil, sentinel
	}
	neverCalled := func(entries []PortEntry) ([]PortEntry, error) {
		called = true
		return entries, nil
	}
	p, _ := NewPipeline(failStep, neverCalled)
	_, err := p.Run(pipelineEntries())
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if called {
		t.Fatal("subsequent step must not be called after error")
	}
}

func TestPipelineRunEmptyInput(t *testing.T) {
	p, _ := NewPipeline(func(e []PortEntry) ([]PortEntry, error) { return e, nil })
	result, err := p.Run([]PortEntry{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}
