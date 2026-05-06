package ports

import (
	"testing"
)

func inspectorSampleEntries() []PortEntry {
	return []PortEntry{
		{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
		{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "apache"},
		{Port: 443, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
		{Port: 9999, Protocol: "udp", State: "LISTEN", Process: ""},
	}
}

func TestNewInspectorValid(t *testing.T) {
	insp, err := NewInspector(nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if insp == nil {
		t.Fatal("expected non-nil inspector")
	}
}

func TestInspectServiceName(t *testing.T) {
	r, _ := NewResolver(nil)
	insp, _ := NewInspector(r, nil, nil, nil)
	entries := inspectorSampleEntries()
	res := insp.Inspect(entries[0], entries)
	if res.ServiceName == "" {
		t.Error("expected a service name for port 80")
	}
}

func TestInspectNoteSharedPort(t *testing.T) {
	insp, _ := NewInspector(nil, nil, nil, nil)
	entries := inspectorSampleEntries()
	res := insp.Inspect(entries[0], entries)
	if res.Note == "" {
		t.Error("expected a note about shared port 80")
	}
}

func TestInspectNoteNoProcess(t *testing.T) {
	insp, _ := NewInspector(nil, nil, nil, nil)
	entries := inspectorSampleEntries()
	// port 9999 has no process and is unique
	res := insp.Inspect(entries[3], entries)
	if res.Note != "no owning process detected" {
		t.Errorf("unexpected note: %q", res.Note)
	}
}

func TestInspectNoteClean(t *testing.T) {
	insp, _ := NewInspector(nil, nil, nil, nil)
	entries := inspectorSampleEntries()
	// port 443 is unique and has a process
	res := insp.Inspect(entries[2], entries)
	if res.Note != "" {
		t.Errorf("expected empty note, got %q", res.Note)
	}
}

func TestInspectScore(t *testing.T) {
	s, _ := NewScorer(DefaultScoreWeights())
	insp, _ := NewInspector(nil, s, nil, nil)
	entries := inspectorSampleEntries()
	res := insp.Inspect(entries[0], entries)
	if res.Score <= 0 {
		t.Errorf("expected positive score, got %f", res.Score)
	}
}

func TestInspectEntryPreserved(t *testing.T) {
	insp, _ := NewInspector(nil, nil, nil, nil)
	entries := inspectorSampleEntries()
	res := insp.Inspect(entries[2], entries)
	if res.Entry.Port != 443 {
		t.Errorf("expected port 443, got %d", res.Entry.Port)
	}
}
