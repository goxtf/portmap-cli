package ports_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/portmap-cli/internal/ports"
)

func integrationSnapshotEntries() []ports.PortEntry {
	return []ports.PortEntry{
		{Protocol: "tcp", LocalAddress: "0.0.0.0", Port: 22, Process: "sshd", State: "LISTEN"},
		{Protocol: "tcp", LocalAddress: "0.0.0.0", Port: 3000, Process: "node", State: "LISTEN"},
		{Protocol: "udp", LocalAddress: "0.0.0.0", Port: 123, Process: "ntpd", State: ""},
	}
}

func TestSnapshotFileIsValidJSON(t *testing.T) {
	dir := t.TempDir()
	sm, err := ports.NewSnapshotManager(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, err := sm.Save(integrationSnapshotEntries())
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	var snap ports.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestSnapshotFilenameContainsTimestamp(t *testing.T) {
	dir := t.TempDir()
	sm, _ := ports.NewSnapshotManager(dir)
	path, _ := sm.Save(integrationSnapshotEntries())
	base := filepath.Base(path)
	if !strings.HasPrefix(base, "snapshot_") {
		t.Errorf("expected filename to start with 'snapshot_', got: %s", base)
	}
}

func TestSnapshotRoundTrip(t *testing.T) {
	dir := t.TempDir()
	sm, _ := ports.NewSnapshotManager(dir)
	original := integrationSnapshotEntries()
	path, _ := sm.Save(original)
	snap, err := sm.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(snap.Entries) != len(original) {
		t.Fatalf("entry count mismatch: got %d, want %d", len(snap.Entries), len(original))
	}
	for i, e := range snap.Entries {
		if e.Port != original[i].Port || e.Process != original[i].Process {
			t.Errorf("entry %d mismatch", i)
		}
	}
}

func TestSnapshotTimestampIsUTC(t *testing.T) {
	dir := t.TempDir()
	sm, _ := ports.NewSnapshotManager(dir)
	path, _ := sm.Save(integrationSnapshotEntries())
	snap, _ := sm.Load(path)
	if snap.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got: %v", snap.Timestamp.Location())
	}
}
