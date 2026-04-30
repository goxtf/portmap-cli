package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestProfileCommandHelp(t *testing.T) {
	root := rootCmd
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"profile", "--help"})
	_ = root.Execute()
	out := buf.String()
	if !strings.Contains(out, "profile") {
		t.Errorf("expected 'profile' in help output, got: %s", out)
	}
}

func TestProfileCmdRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "profile" {
			found = true
			break
		}
	}
	if !found {
		t.Error("profile command not registered on root")
	}
}

func TestProfileTopFlag(t *testing.T) {
	cmd, _, _ := rootCmd.Find([]string{"profile"})
	if cmd == nil {
		t.Fatal("profile command not found")
	}
	f := cmd.Flags().Lookup("top")
	if f == nil {
		t.Fatal("expected --top flag")
	}
	if f.DefValue != "5" {
		t.Errorf("expected default 5, got %s", f.DefValue)
	}
}

func TestProfileJSONFlag(t *testing.T) {
	cmd, _, _ := rootCmd.Find([]string{"profile"})
	if cmd == nil {
		t.Fatal("profile command not found")
	}
	f := cmd.Flags().Lookup("json")
	if f == nil {
		t.Fatal("expected --json flag")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}
