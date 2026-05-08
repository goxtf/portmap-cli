package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findValidateCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Use == "validate" {
			return c
		}
	}
	return nil
}

func TestValidateCommandHelp(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"validate", "--help"})
	_ = rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "validate") {
		t.Errorf("expected 'validate' in help output, got: %s", out)
	}
}

func TestValidateCmdRegistered(t *testing.T) {
	if findValidateCmd(rootCmd) == nil {
		t.Fatal("validate command not registered on rootCmd")
	}
}

func TestValidateProtocolsFlag(t *testing.T) {
	cmd := findValidateCmd(rootCmd)
	if cmd == nil {
		t.Fatal("validate command not found")
	}
	f := cmd.Flags().Lookup("protocols")
	if f == nil {
		t.Fatal("--protocols flag not found")
	}
	if f.DefValue != "" {
		t.Errorf("expected empty default for --protocols, got %q", f.DefValue)
	}
}

func TestValidateMaxPortFlag(t *testing.T) {
	cmd := findValidateCmd(rootCmd)
	if cmd == nil {
		t.Fatal("validate command not found")
	}
	f := cmd.Flags().Lookup("max-port")
	if f == nil {
		t.Fatal("--max-port flag not found")
	}
	if f.DefValue != "65535" {
		t.Errorf("expected default 65535 for --max-port, got %q", f.DefValue)
	}
}

func TestValidateRequireProcessFlag(t *testing.T) {
	cmd := findValidateCmd(rootCmd)
	if cmd == nil {
		t.Fatal("validate command not found")
	}
	f := cmd.Flags().Lookup("require-process")
	if f == nil {
		t.Fatal("--require-process flag not found")
	}
}

func TestValidateJSONFlag(t *testing.T) {
	cmd := findValidateCmd(rootCmd)
	if cmd == nil {
		t.Fatal("validate command not found")
	}
	f := cmd.Flags().Lookup("json")
	if f == nil {
		t.Fatal("--json flag not found")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false for --json, got %q", f.DefValue)
	}
}
