package ports

import (
	"testing"
)

var taggerSampleEntries = []PortEntry{
	{Protocol: "tcp", Port: 80, State: "LISTEN", Process: "nginx"},
	{Protocol: "udp", Port: 53, State: "UNCONN", Process: "dnsmasq"},
	{Protocol: "tcp", Port: 443, State: "LISTEN", Process: "nginx"},
	{Protocol: "tcp", Port: 8080, State: "ESTABLISHED", Process: "java"},
}

func TestNewTaggerValidRules(t *testing.T) {
	rules := []TagRule{{Field: "protocol", Match: "tcp", Key: "proto-tag", Value: "tcp"}}
	_, err := NewTagger(rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewTaggerInvalidField(t *testing.T) {
	rules := []TagRule{{Field: "invalid", Match: "tcp", Key: "k", Value: "v"}}
	_, err := NewTagger(rules)
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestNewTaggerEmptyKey(t *testing.T) {
	rules := []TagRule{{Field: "protocol", Match: "tcp", Key: "", Value: "v"}}
	_, err := NewTagger(rules)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTagByProtocol(t *testing.T) {
	rules := []TagRule{{Field: "protocol", Match: "tcp", Key: "type", Value: "tcp"}}
	tagger, _ := NewTagger(rules)
	result := tagger.Tag(taggerSampleEntries)

	// entries 0, 2, 3 are tcp
	for _, idx := range []int{0, 2, 3} {
		tags := result[idx]
		if len(tags) == 0 || tags[0].Value != "tcp" {
			t.Errorf("entry %d expected tag tcp", idx)
		}
	}
	if len(result[1]) != 0 {
		t.Error("udp entry should not be tagged")
	}
}

func TestTagByProcess(t *testing.T) {
	rules := []TagRule{{Field: "process", Match: "nginx", Key: "webserver", Value: "true"}}
	tagger, _ := NewTagger(rules)
	result := tagger.Tag(taggerSampleEntries)

	if len(result[0]) == 0 || result[0][0].Key != "webserver" {
		t.Error("expected nginx entry to be tagged as webserver")
	}
	if len(result[1]) != 0 {
		t.Error("dnsmasq should not be tagged as webserver")
	}
}

func TestTagMultipleRules(t *testing.T) {
	rules := []TagRule{
		{Field: "protocol", Match: "tcp", Key: "proto", Value: "tcp"},
		{Field: "state", Match: "LISTEN", Key: "listening", Value: "yes"},
	}
	tagger, _ := NewTagger(rules)
	result := tagger.Tag(taggerSampleEntries)

	// entry 0: tcp + LISTEN → 2 tags
	if len(result[0]) != 2 {
		t.Errorf("expected 2 tags for entry 0, got %d", len(result[0]))
	}
}

func TestTagByPort(t *testing.T) {
	rules := []TagRule{{Field: "port", Match: "53", Key: "dns", Value: "true"}}
	tagger, _ := NewTagger(rules)
	result := tagger.Tag(taggerSampleEntries)

	if len(result[1]) == 0 || result[1][0].Key != "dns" {
		t.Error("expected port 53 entry to be tagged as dns")
	}
}
