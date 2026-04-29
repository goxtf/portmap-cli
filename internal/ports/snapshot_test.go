package ports

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func snapshotEntries() []PortEntry {
	return []PortEntry{
		{Protocol: "tcp", LocalAddress: "0.0.0.0", Port: 80, Process: "nginx", State: "LISTEN"},
		{Protocol: "tcp", LocalAddress: "0.0.0.0", Port: 443, Process: "nginx", State: "LISTEN"},
		{Protocol: "udp", LocalAddress: "127.0.0.1", Port: 53, Process: "dns", State: ""},
	}
}

func TestNewSnapshotManagerEmptyDir(t *testing.T) {
	_, err := NewSnapshotManager("")
	if err == nil {
		t.Fatal("expected error for empty dir")
	}
}

func TestNewSnapshotManagerCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snapshots")
	sm, err := NewSnapshotManager(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sm == nil {
		t.Fatal("expected non-nil manager")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	sm, _ := NewSnapshotManager(dir)
	entries := snapshotEntries()

	path, err := sm.Save(entries)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	snap, err := sm.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(snap.Entries) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(snap.Entries))
	}
	if snap.Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestSnapshotDiff(t *testing.T) {
	dir := t.TempDir()
	sm, _ := NewSnapshotManager(dir)

	old := &Snapshot{Timestamp: time.Now(), Entries: snapshotEntries()}
	newEntries := snapshotEntries()
	newEntries = append(newEntries, PortEntry{Protocol: "tcp", Port: 8080, Process: "app", State: "LISTEN"})
	newEntries = newEntries[1:] // remove first entry
	newSnap := &Snapshot{Timestamp: time.Now(), Entries: newEntries}

	added, removed := sm.Diff(old, newSnap)
	if len(added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(added))
	}
	if len(removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(removed))
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	dir := t.TempDir()
	sm, _ := NewSnapshotManager(dir)
	_, err := sm.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error loading nonexistent file")
	}
}
