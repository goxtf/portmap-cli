package ports

import (
	"testing"
)

var labelerEntries = []PortEntry{
	{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
	{Protocol: "udp", Port: 53, State: "UNCONN", Process: "dnsmasq"},
	{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
	{Protocol: "tcp", Port: 8080, State: "LISTEN", Process: "java"},
}

func TestNewLabelerValidRules(t *testing.T) {
	rules := []LabelRule{{Field: "protocol", Match: "tcp", Label: "tcp-traffic"}}
	_, err := NewLabeler(rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewLabelerInvalidField(t *testing.T) {
	rules := []LabelRule{{Field: "hostname", Match: "localhost", Label: "local"}}
	_, err := NewLabeler(rules)
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestNewLabelerEmptyMatch(t *testing.T) {
	rules := []LabelRule{{Field: "protocol", Match: "", Label: "any"}}
	_, err := NewLabeler(rules)
	if err == nil {
		t.Fatal("expected error for empty match")
	}
}

func TestNewLabelerEmptyLabel(t *testing.T) {
	rules := []LabelRule{{Field: "protocol", Match: "tcp", Label: ""}}
	_, err := NewLabeler(rules)
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestLabelByProtocol(t *testing.T) {
	rules := []LabelRule{{Field: "protocol", Match: "tcp", Label: "tcp-traffic"}}
	lb, _ := NewLabeler(rules)
	out := lb.Apply(labelerEntries)
	for _, e := range out {
		if e.Protocol == "tcp" {
			if _, ok := e.Tags["label:tcp-traffic"]; !ok {
				t.Errorf("expected label tcp-traffic on tcp entry port %d", e.Port)
			}
		} else {
			if _, ok := e.Tags["label:tcp-traffic"]; ok {
				t.Errorf("unexpected label tcp-traffic on non-tcp entry port %d", e.Port)
			}
		}
	}
}

func TestLabelByProcess(t *testing.T) {
	rules := []LabelRule{{Field: "process", Match: "nginx", Label: "web"}}
	lb, _ := NewLabeler(rules)
	out := lb.Apply(labelerEntries)
	webCount := 0
	for _, e := range out {
		if _, ok := e.Tags["label:web"]; ok {
			webCount++
		}
	}
	if webCount != 2 {
		t.Errorf("expected 2 web-labeled entries, got %d", webCount)
	}
}

func TestLabelMultipleRulesOnSameEntry(t *testing.T) {
	rules := []LabelRule{
		{Field: "protocol", Match: "tcp", Label: "tcp-traffic"},
		{Field: "state", Match: "LISTEN", Label: "listening"},
	}
	lb, _ := NewLabeler(rules)
	out := lb.Apply(labelerEntries)
	for _, e := range out {
		if e.Protocol == "tcp" && e.State == "LISTEN" {
			if _, ok := e.Tags["label:tcp-traffic"]; !ok {
				t.Errorf("missing tcp-traffic label on entry port %d", e.Port)
			}
			if _, ok := e.Tags["label:listening"]; !ok {
				t.Errorf("missing listening label on entry port %d", e.Port)
			}
		}
	}
}

func TestLabelNoRulesNoTags(t *testing.T) {
	lb, _ := NewLabeler([]LabelRule{})
	out := lb.Apply(labelerEntries)
	for _, e := range out {
		if len(e.Tags) != 0 {
			t.Errorf("expected no tags, got %v", e.Tags)
		}
	}
}
