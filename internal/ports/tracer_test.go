package ports

import (
	"testing"
)

var tracerEntries = []PortEntry{
	{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "nginx", Tags: map[string]string{"env": "prod"}},
	{Port: 443, Protocol: "tcp", State: "LISTEN", Process: "nginx", Tags: map[string]string{}},
	{Port: 8080, Protocol: "tcp", State: "CLOSE_WAIT", Process: "", Tags: map[string]string{}},
	{Port: 80, Protocol: "udp", State: "UNCONN", Process: "systemd", Tags: map[string]string{}},
}

func TestNewTracerNilResolver(t *testing.T) {
	tr := NewTracer(nil)
	if tr == nil {
		t.Fatal("expected non-nil Tracer")
	}
}

func TestTracerInvalidPort(t *testing.T) {
	tr := NewTracer(nil)
	_, err := tr.Trace(tracerEntries, 0)
	if err == nil {
		t.Fatal("expected error for port 0")
	}
	_, err = tr.Trace(tracerEntries, 70000)
	if err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestTracerNoMatch(t *testing.T) {
	tr := NewTracer(nil)
	_, err := tr.Trace(tracerEntries, 9999)
	if err == nil {
		t.Fatal("expected error when no entries match")
	}
}

func TestTracerMultipleMatches(t *testing.T) {
	tr := NewTracer(nil)
	results, err := tr.Trace(tracerEntries, 80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestTracerHopsContainPort(t *testing.T) {
	tr := NewTracer(nil)
	results, err := tr.Trace(tracerEntries, 443)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, h := range results[0].Hops {
		if h == "port=443 matched" {
			found = true
		}
	}
	if !found {
		t.Error("expected hop 'port=443 matched' not found")
	}
}

func TestTracerTagsAppearInHops(t *testing.T) {
	tr := NewTracer(nil)
	results, err := tr.Trace(tracerEntries, 80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// first result is the tcp/nginx entry with tag env=prod
	found := false
	for _, r := range results {
		for _, h := range r.Hops {
			if h == "tag:env=prod" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected tag hop 'tag:env=prod' not found")
	}
}

func TestTracerCloseWaitNote(t *testing.T) {
	tr := NewTracer(nil)
	results, err := tr.Trace(tracerEntries, 8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results[0].Notes) == 0 {
		t.Error("expected notes for CLOSE_WAIT state")
	}
}

func TestTracerWithResolver(t *testing.T) {
	r, _ := NewResolver(nil)
	tr := NewTracer(r)
	results, err := tr.Trace(tracerEntries, 80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, res := range results {
		for _, h := range res.Hops {
			if h == "service=http" {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected service=http hop from resolver")
	}
}
