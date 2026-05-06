package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findInspectCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Use == "inspect <port>" {
			return c
		}
	}
	return nil
}

func TestInspectCommandHelp(t *testing.T) {
	cmd := findInspectCmd(rootCmd)
	if cmd == nil {
		t.Fatal("inspect command not registered")
	}
	if !strings.Contains(cmd.Short, "inspect") {
		t.Errorf("unexpected short description: %s", cmd.Short)
	}
}

func TestInspectCmdRegistered(t *testing.T) {
	if findInspectCmd(rootCmd) == nil {
		t.Fatal("inspect command should be registered on rootCmd")
	}
}

func TestInspectJSONFlag(t *testing.T) {
	cmd := findInspectCmd(rootCmd)
	if cmd == nil {
		t.Fatal("inspect command not found")
	}
	f := cmd.Flags().Lookup("json")
	if f == nil {
		t.Fatal("expected --json flag")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}

func TestInspectRequiresPortArg(t *testing.T) {
	cmd := findInspectCmd(rootCmd)
	if cmd == nil {
		t.Fatal("inspect command not found")
	}
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error when no port argument provided")
	}
}

func TestInspectAcceptsOneArg(t *testing.T) {
	cmd := findInspectCmd(rootCmd)
	if cmd == nil {
		t.Fatal("inspect command not found")
	}
	err := cmd.Args(cmd, []string{"8080"})
	if err != nil {
		t.Errorf("unexpected error for valid arg: %v", err)
	}
}
