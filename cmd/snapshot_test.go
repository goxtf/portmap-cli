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

func TestSnapshotDirFlag(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"snapshot", "--help"})
	_ = cmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "--dir") {
		t.Errorf("expected --dir flag in help, got: %s", out)
	}
}

func TestSnapshotDiffFlag(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"snapshot", "--help"})
	_ = cmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "--diff") {
		t.Errorf("expected --diff flag in help, got: %s", out)
	}
}
