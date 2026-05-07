package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestSnapshotCommandHelp(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"snapshot", "--help"})
	_ = cmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "snapshot") {
		t.Errorf("expected 'snapshot' in help output, got: %s", out)
	}
}

func TestSnapshotCmdRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "snapshot" {
			return
		}
	}
	t.Fatal("snapshot command not registered")
}

// snapshotHelpOutput is a helper that runs the snapshot --help command and
// returns the output string, reducing duplication across flag tests.
func snapshotHelpOutput(t *testing.T) string {
	t.Helper()
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"snapshot", "--help"})
	_ = cmd.Execute()
	return buf.String()
}

func TestSnapshotDirFlag(t *testing.T) {
	out := snapshotHelpOutput(t)
	if !strings.Contains(out, "--dir") {
		t.Errorf("expected --dir flag in help, got: %s", out)
	}
}

func TestSnapshotDiffFlag(t *testing.T) {
	out := snapshotHelpOutput(t)
	if !strings.Contains(out, "--diff") {
		t.Errorf("expected --diff flag in help, got: %s", out)
	}
}
