package ports

import (
	"testing"
)

var classifierSampleEntries = []PortEntry{
	{Protocol: "tcp", State: "LISTEN", Process: "nginx", Port: 80},
	{Protocol: "udp", State: "UNCONN", Process: "systemd-resolved", Port: 53},
	{Protocol: "tcp", State: "ESTABLISHED", Process: "postgres", Port: 5432},
	{Protocol: "tcp", State: "LISTEN", Process: "sshd", Port: 22},
}

func TestNewClassifierValidRules(t *testing.T) {
	rules := []ClassifierRule{
		{Field: "protocol", Contains: "udp", Category: "dns"},
	}
	_, err := NewClassifier(rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewClassifierInvalidField(t *testing.T) {
	rules := []ClassifierRule{
		{Field: "hostname", Contains: "foo", Category: "bar"},
	}
	_, err := NewClassifier(rules)
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestNewClassifierEmptyCategory(t *testing.T) {
	rules := []ClassifierRule{
		{Field: "protocol", Contains: "tcp", Category: ""},
	}
	_, err := NewClassifier(rules)
	if err == nil {
		t.Fatal("expected error for empty category")
	}
}

func TestClassifyByProtocol(t *testing.T) {
	rules := []ClassifierRule{
		{Field: "protocol", Contains: "udp", Category: "udp-service"},
		{Field: "protocol", Contains: "tcp", Category: "tcp-service"},
	}
	c, _ := NewClassifier(rules)
	out := c.Classify(classifierSampleEntries)

	if out[0].Tags["category"] != "tcp-service" {
		t.Errorf("expected tcp-service, got %s", out[0].Tags["category"])
	}
	if out[1].Tags["category"] != "udp-service" {
		t.Errorf("expected udp-service, got %s", out[1].Tags["category"])
	}
}

func TestClassifyByProcess(t *testing.T) {
	rules := []ClassifierRule{
		{Field: "process", Contains: "nginx", Category: "web"},
		{Field: "process", Contains: "postgres", Category: "database"},
	}
	c, _ := NewClassifier(rules)
	out := c.Classify(classifierSampleEntries)

	if out[0].Tags["category"] != "web" {
		t.Errorf("expected web, got %s", out[0].Tags["category"])
	}
	if out[2].Tags["category"] != "database" {
		t.Errorf("expected database, got %s", out[2].Tags["category"])
	}
}

func TestClassifyFallsBackToOther(t *testing.T) {
	rules := []ClassifierRule{
		{Field: "process", Contains: "nginx", Category: "web"},
	}
	c, _ := NewClassifier(rules)
	out := c.Classify(classifierSampleEntries)

	if out[3].Tags["category"] != "other" {
		t.Errorf("expected other for sshd, got %s", out[3].Tags["category"])
	}
}

func TestClassifyDoesNotMutateOriginal(t *testing.T) {
	rules := []ClassifierRule{
		{Field: "state", Contains: "listen", Category: "listening"},
	}
	c, _ := NewClassifier(rules)
	original := make([]PortEntry, len(classifierSampleEntries))
	copy(original, classifierSampleEntries)
	c.Classify(classifierSampleEntries)
	for i, e := range classifierSampleEntries {
		if e.Tags != nil && original[i].Tags == nil {
			t.Errorf("original entry %d was mutated", i)
		}
	}
}
