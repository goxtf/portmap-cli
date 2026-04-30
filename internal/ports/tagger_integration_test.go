package ports

import "testing"

func TestTaggerNoRulesProducesEmptyResult(t *testing.T) {
	tagger, err := NewTagger([]TagRule{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := []PortEntry{
		{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
	}
	result := tagger.Tag(entries)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries tagged", len(result))
	}
}

func TestTaggerAllEntriesTagged(t *testing.T) {
	entries := []PortEntry{
		{Protocol: "tcp", Port: 22, State: "LISTEN", Process: "sshd"},
		{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
		{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
	}
	rules := []TagRule{
		{Field: "protocol", Match: "tcp", Key: "proto", Value: "tcp"},
	}
	tagger, _ := NewTagger(rules)
	result := tagger.Tag(entries)

	if len(result) != len(entries) {
		t.Errorf("expected all %d entries tagged, got %d", len(entries), len(result))
	}
}

func TestTaggerTagValuesAreCorrect(t *testing.T) {
	entries := []PortEntry{
		{Protocol: "udp", Port: 123, State: "UNCONN", Process: "ntpd"},
	}
	rules := []TagRule{
		{Field: "process", Match: "ntpd", Key: "service", Value: "ntp"},
	}
	tagger, _ := NewTagger(rules)
	result := tagger.Tag(entries)

	if len(result[0]) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(result[0]))
	}
	if result[0][0].Key != "service" || result[0][0].Value != "ntp" {
		t.Errorf("unexpected tag: %+v", result[0][0])
	}
}
