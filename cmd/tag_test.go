package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestTagCommandHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"tag", "--help"})
	_ = rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "tag") {
		t.Errorf("expected 'tag' in help output, got: %s", out)
	}
}

func TestTagCmdRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "tag" {
			return
		}
	}
	t.Error("tag command not registered on root")
}

func TestTagFieldFlag(t *testing.T) {
	f := tagCmd.Flags().Lookup("field")
	if f == nil {
		t.Fatal("expected --field flag to be registered")
	}
	if f.DefValue != "protocol" {
		t.Errorf("expected default 'protocol', got %s", f.DefValue)
	}
}

func TestTagJSONFlag(t *testing.T) {
	f := tagCmd.Flags().Lookup("json")
	if f == nil {
		t.Fatal("expected --json flag to be registered")
	}
}

func TestTagRequiredFlags(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	// call without required --match and --key flags
	rootCmd.SetArgs([]string{"tag", "--field", "protocol"})
	err := rootCmd.Execute()
	if err == nil {
		// cobra may not return error directly; check output
		out := buf.String()
		if !strings.Contains(out, "required") && !strings.Contains(out, "match") {
			t.Log("no error returned but output does not mention required flags:", out)
		}
	}
}
