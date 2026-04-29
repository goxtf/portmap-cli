package ports

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of port entries.
type Snapshot struct {
	Timestamp time.Time   `json:"timestamp"`
	Entries   []PortEntry `json:"entries"`
}

// SnapshotManager handles saving and loading port snapshots.
type SnapshotManager struct {
	dir string
}

// NewSnapshotManager creates a SnapshotManager that stores snapshots in dir.
func NewSnapshotManager(dir string) (*SnapshotManager, error) {
	if dir == "" {
		return nil, fmt.Errorf("snapshot directory must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create snapshot dir: %w", err)
	}
	return &SnapshotManager{dir: dir}, nil
}

// Save writes a snapshot of entries to a timestamped JSON file.
func (sm *SnapshotManager) Save(entries []PortEntry) (string, error) {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Entries:   entries,
	}
	filename := fmt.Sprintf("%s/snapshot_%s.json", sm.dir, snap.Timestamp.Format("20060102_150405"))
	f, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("create snapshot file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return "", fmt.Errorf("encode snapshot: %w", err)
	}
	return filename, nil
}

// Load reads a snapshot from the given file path.
func (sm *SnapshotManager) Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open snapshot file: %w", err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("decode snapshot: %w", err)
	}
	return &snap, nil
}

// Diff returns added and removed entries between two snapshots.
func (sm *SnapshotManager) Diff(old, new *Snapshot) (added, removed []PortEntry) {
	return diffEntries(old.Entries, new.Entries)
}
