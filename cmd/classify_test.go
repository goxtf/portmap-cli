package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findClassifyCmd(root *cobra.Command) *cobra.Command {
	for _, sub := range root.Commands() {
		if sub.Use == "classify" {
			return sub
		}
	}
	return nil
}

func TestClassifyCommandHelp(t *testing.T) {
	root := buildRootCmd()
	root.SetArgs([]string{"classify", "--help"})
	var sb strings.Builder
	root.SetOut(&sb)
	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "classify") {
		t.Error("help output should mention 'classify'")
	}
}

func TestClassifyCmdRegistered(t *testing.T) {
	root := buildRootCmd()
	cmd := findClassifyCmd(root)
	if cmd == nil {
		t.Fatal("classify command not registered")
	}
}

func TestClassifyFieldFlag(t *testing.T) {
	root := buildRootCmd()
	cmd := findClassifyCmd(root)
	if cmd == nil {
		t.Fatal("classify command not found")
	}
	f := cmd.Flags().Lookup("field")
	if f == nil {
		t.Fatal("--field flag not registered")
	}
	if f.DefValue != "process" {
		t.Errorf("expected default 'process', got %s", f.DefValue)
	}
}

func TestClassifyRequiredFlags(t *testing.T) {
	root := buildRootCmd()
	root.SetArgs([]string{"classify", "--field", "protocol"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when required flags are missing")
	}
}

func TestClassifyJSONFlag(t *testing.T) {
	root := buildRootCmd()
	cmd := findClassifyCmd(root)
	if cmd == nil {
		t.Fatal("classify command not found")
	}
	f := cmd.Flags().Lookup("json")
	if f == nil {
		t.Fatal("--json flag not registered")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}
