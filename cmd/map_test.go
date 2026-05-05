package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findMapCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Use == "map" {
			return c
		}
	}
	return nil
}

func TestMapCommandHelp(t *testing.T) {
	cmd := findMapCmd(rootCmd)
	if cmd == nil {
		t.Fatal("map command not registered")
	}
	if !strings.Contains(cmd.Short, "port") {
		t.Errorf("expected Short to mention 'port', got: %s", cmd.Short)
	}
}

func TestMapCmdRegistered(t *testing.T) {
	if findMapCmd(rootCmd) == nil {
		t.Fatal("expected map command to be registered on rootCmd")
	}
}

func TestMapIncludeClosedFlag(t *testing.T) {
	cmd := findMapCmd(rootCmd)
	if cmd == nil {
		t.Fatal("map command not found")
	}
	f := cmd.Flags().Lookup("include-closed")
	if f == nil {
		t.Fatal("expected --include-closed flag to exist")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}

func TestMapConflictsFlag(t *testing.T) {
	cmd := findMapCmd(rootCmd)
	if cmd == nil {
		t.Fatal("map command not found")
	}
	f := cmd.Flags().Lookup("conflicts")
	if f == nil {
		t.Fatal("expected --conflicts flag to exist")
	}
}

func TestMapPortFlag(t *testing.T) {
	cmd := findMapCmd(rootCmd)
	if cmd == nil {
		t.Fatal("map command not found")
	}
	f := cmd.Flags().Lookup("port")
	if f == nil {
		t.Fatal("expected --port flag to exist")
	}
	if f.DefValue != "0" {
		t.Errorf("expected default 0, got %s", f.DefValue)
	}
}
