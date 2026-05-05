package ports

import "testing"

func TestLabelerPreservesExistingTags(t *testing.T) {
	entries := []PortEntry{
		{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx", Tags: map[string]string{"env": "prod"}},
	}
	rules := []LabelRule{{Field: "protocol", Match: "tcp", Label: "tcp-traffic"}}
	lb, _ := NewLabeler(rules)
	out := lb.Apply(entries)
	if out[0].Tags["env"] != "prod" {
		t.Errorf("expected existing tag env=prod to be preserved")
	}
	if _, ok := out[0].Tags["label:tcp-traffic"]; !ok {
		t.Errorf("expected label:tcp-traffic to be added")
	}
}

func TestLabelerDoesNotMutateOriginal(t *testing.T) {
	original := []PortEntry{
		{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
	}
	rules := []LabelRule{{Field: "protocol", Match: "tcp", Label: "secure"}}
	lb, _ := NewLabeler(rules)
	lb.Apply(original)
	if original[0].Tags != nil {
		t.Errorf("original entry should not be mutated")
	}
}

func TestLabelerPortFieldMatch(t *testing.T) {
	entries := []PortEntry{
		{Protocol: "tcp", Port: 22, State: "LISTEN", Process: "sshd"},
		{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
	}
	rules := []LabelRule{{Field: "port", Match: "22", Label: "ssh"}}
	lb, _ := NewLabeler(rules)
	out := lb.Apply(entries)
	if _, ok := out[0].Tags["label:ssh"]; !ok {
		t.Errorf("expected ssh label on port 22")
	}
	if _, ok := out[1].Tags["label:ssh"]; ok {
		t.Errorf("unexpected ssh label on port 80")
	}
}
