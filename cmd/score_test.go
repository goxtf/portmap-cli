package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findScoreCmd(root *cobra.Command) *cobra.Command {
	for _, sub := range root.Commands() {
		if sub.Use == "score" {
			return sub
		}
	}
	return nil
}

func TestScoreCommandHelp(t *testing.T) {
	cmd := findScoreCmd(rootCmd)
	if cmd == nil {
		t.Fatal("score command not registered")
	}
	if !strings.Contains(cmd.Short, "Rank") {
		t.Errorf("unexpected Short description: %q", cmd.Short)
	}
}

func TestScoreCmdRegistered(t *testing.T) {
	if findScoreCmd(rootCmd) == nil {
		t.Fatal("expected score command to be registered on rootCmd")
	}
}

func TestScoreWeightPrivFlag(t *testing.T) {
	cmd := findScoreCmd(rootCmd)
	if cmd == nil {
		t.Fatal("score command not found")
	}
	f := cmd.Flags().Lookup("weight-priv")
	if f == nil {
		t.Fatal("expected --weight-priv flag")
	}
	if f.DefValue != "3" {
		t.Errorf("expected default 3, got %s", f.DefValue)
	}
}

func TestScoreTopFlag(t *testing.T) {
	cmd := findScoreCmd(rootCmd)
	if cmd == nil {
		t.Fatal("score command not found")
	}
	f := cmd.Flags().Lookup("top")
	if f == nil {
		t.Fatal("expected --top flag")
	}
	if f.DefValue != "0" {
		t.Errorf("expected default 0, got %s", f.DefValue)
	}
}

func TestScoreJSONFlag(t *testing.T) {
	cmd := findScoreCmd(rootCmd)
	if cmd == nil {
		t.Fatal("score command not found")
	}
	f := cmd.Flags().Lookup("json")
	if f == nil {
		t.Fatal("expected --json flag")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}
