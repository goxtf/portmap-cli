package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// resetListFlags resets the list command flags to their default values
// between tests to avoid state leakage.
func resetListFlags() {
	filterProtocol = ""
	filterState = ""
	filterProcess = ""
	filterPort = 0
	outputFormat = "table"
}

// findListCmd is a helper that retrieves the "list" subcommand from rootCmd.
// It returns nil if the command is not registered.
func findListCmd() *cobra.Command {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "list" {
			return cmd
		}
	}
	return nil
}

func TestListCommandHelp(t *testing.T) {
	defer resetListFlags()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list", "--help"})

	// Execute should not return an error for --help
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error for --help, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "list") {
		t.Errorf("expected help output to contain 'list', got: %s", output)
	}
}

func TestListCmdRegistered(t *testing.T) {
	var found bool
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "list" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'list' command to be registered on root")
	}
}

func TestListProtocolFlag(t *testing.T) {
	defer resetListFlags()

	listCmd := findListCmd()
	if listCmd == nil {
		t.Fatal("list command not found")
	}

	f := listCmd.Flags().Lookup("protocol")
	if f == nil {
		t.Fatal("expected --protocol flag to exist on list command")
	}
	if f.DefValue != "" {
		t.Errorf("expected default value of --protocol to be empty, got: %s", f.DefValue)
	}
}

func TestListStateFlag(t *testing.T) {
	defer resetListFlags()

	listCmd := findListCmd()
	if listCmd == nil {
		t.Fatal("list command not found")
	}

	f := listCmd.Flags().Lookup("state")
	if f == nil {
		t.Fatal("expected --state flag to exist on list command")
	}
	if f.DefValue != "" {
		t.Errorf("expected default value of --state to be empty, got: %s", f.DefValue)
	}
}

func TestListFormatFlag(t *testing.T) {
	defer resetListFlags()

	listCmd := findListCmd()
	if listCmd == nil {
		t.Fatal("list command not found")
	}

	f := listCmd.Flags().Lookup("format")
	if f == nil {
		t.Fatal("expected --format flag to exist on list command")
	}
	if f.DefValue != "table" {
		t.Errorf("expected default format to be 'table', got: %s", f.DefValue)
	}
}

func TestListPortFlag(t *testing.T) {
	defer resetListFlags()

	listCmd := findListCmd()
	if listCmd == nil {
		t.Fatal("list command not found")
	}

	f := listCmd.Flags().Lookup("port")
	if f == nil {
		t.Fatal("expected --port flag to exist on list command")
	}
	if f.DefValue != "0" {
		t.Errorf("expected default port to be '0', got: %s", f.DefValue)
	}
}
