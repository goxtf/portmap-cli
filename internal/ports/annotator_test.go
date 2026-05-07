package ports

import (
	"testing"
)

var annotatorSampleEntries = []PortEntry{
	{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 443, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
	{Port: 5432, Protocol: "tcp", State: "LISTEN", Process: "postgres"},
	{Port: 8080, Protocol: "tcp", State: "CLOSE_WAIT", Process: "java"},
}

func TestNewAnnotatorNoRules(t *testing.T) {
	_, err := NewAnnotator(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNewAnnotatorInvalidField(t *testing.T) {
	_, err := NewAnnotator([]AnnotationRule{{Field: "invalid", Match: "tcp", Note: "note"}})
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
}

func TestNewAnnotatorEmptyMatch(t *testing.T) {
	_, err := NewAnnotator([]AnnotationRule{{Field: "protocol", Match: "", Note: "note"}})
	if err == nil {
		t.Fatal("expected error for empty match")
	}
}

func TestNewAnnotatorEmptyNote(t *testing.T) {
	_, err := NewAnnotator([]AnnotationRule{{Field: "protocol", Match: "tcp", Note: ""}})
	if err == nil {
		t.Fatal("expected error for empty note")
	}
}

func TestAnnotateByProtocol(t *testing.T) {
	a, err := NewAnnotator([]AnnotationRule{{Field: "protocol", Match: "tcp", Note: "TCP traffic"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := a.Annotate(annotatorSampleEntries)
	for _, e := range out {
		if e.Protocol == "tcp" && e.Notes != "TCP traffic" {
			t.Errorf("expected note 'TCP traffic', got %q", e.Notes)
		}
	}
}

func TestAnnotateByProcess(t *testing.T) {
	a, _ := NewAnnotator([]AnnotationRule{{Field: "process", Match: "postgres", Note: "DB process"}})
	out := a.Annotate(annotatorSampleEntries)
	for _, e := range out {
		if e.Process == "postgres" && e.Notes != "DB process" {
			t.Errorf("expected note 'DB process', got %q", e.Notes)
		}
		if e.Process != "postgres" && e.Notes != "" {
			t.Errorf("expected no note for non-postgres entry, got %q", e.Notes)
		}
	}
}

func TestAnnotateMultipleRulesAppend(t *testing.T) {
	a, _ := NewAnnotator([]AnnotationRule{
		{Field: "protocol", Match: "tcp", Note: "TCP"},
		{Field: "state", Match: "LISTEN", Note: "listening"},
	})
	out := a.Annotate(annotatorSampleEntries)
	for _, e := range out {
		if e.Protocol == "tcp" && e.State == "LISTEN" {
			if e.Notes != "TCP; listening" {
				t.Errorf("expected 'TCP; listening', got %q", e.Notes)
			}
		}
	}
}

func TestAnnotateDoesNotMutateOriginal(t *testing.T) {
	original := make([]PortEntry, len(annotatorSampleEntries))
	copy(original, annotatorSampleEntries)
	a, _ := NewAnnotator([]AnnotationRule{{Field: "process", Match: "nginx", Note: "web server"}})
	a.Annotate(annotatorSampleEntries)
	for i, e := range annotatorSampleEntries {
		if e.Notes != original[i].Notes {
			t.Errorf("original entry mutated at index %d", i)
		}
	}
}
