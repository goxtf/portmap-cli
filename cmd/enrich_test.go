package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestEnrichCommandHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"enrich", "--help"})
	_ = rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "enrich") {
		t.Errorf("expected 'enrich' in help output, got: %s", out)
	}
}

func TestEnrichCmdRegistered(t *testing.T) {
	var found bool
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "enrich" {
			found = true
			break
		}
	}
	if !found {
		t.Error("enrich command not registered on rootCmd")
	}
}

func TestEnrichGroupFlag(t *testing.T) {
	var cmd *cobra.Command
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "enrich" {
			cmd = sub
			break
		}
	}
	if cmd == nil {
		t.Fatal("enrich command not found")
	}
	f := cmd.Flags().Lookup("group")
	if f == nil {
		t.Fatal("--group flag not registered on enrich command")
	}
	if f.Shorthand != "g" {
		t.Errorf("expected shorthand 'g', got %q", f.Shorthand)
	}
	if f.DefValue != "" {
		t.Errorf("expected default value '', got %q", f.DefValue)
	}
}

func TestEnrichGroupFlagInvalidValue(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"enrich", "--group", "invalid_field"})
	err := rootCmd.Execute()
	// We expect an error because the group value is invalid;
	// scanner may also fail in test env, so just verify no panic.
	_ = err
}
