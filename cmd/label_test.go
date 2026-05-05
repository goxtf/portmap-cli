package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findLabelCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Use == "label" {
			return c
		}
	}
	return nil
}

func TestLabelCommandHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"label", "--help"})
	_ = rootCmd.Execute()
	output := buf.String()
	if !strings.Contains(output, "label") {
		t.Errorf("expected 'label' in help output, got: %s", output)
	}
}

func TestLabelCmdRegistered(t *testing.T) {
	c := findLabelCmd(rootCmd)
	if c == nil {
		t.Fatal("label command not registered")
	}
}

func TestLabelFieldFlag(t *testing.T) {
	c := findLabelCmd(rootCmd)
	if c == nil {
		t.Fatal("label command not registered")
	}
	f := c.Flags().Lookup("field")
	if f == nil {
		t.Fatal("expected --field flag")
	}
	if f.DefValue != "protocol" {
		t.Errorf("expected default field=protocol, got %s", f.DefValue)
	}
}

func TestLabelJSONFlag(t *testing.T) {
	c := findLabelCmd(rootCmd)
	if c == nil {
		t.Fatal("label command not registered")
	}
	f := c.Flags().Lookup("json")
	if f == nil {
		t.Fatal("expected --json flag")
	}
}

func TestLabelRequiredFlags(t *testing.T) {
	c := findLabelCmd(rootCmd)
	if c == nil {
		t.Fatal("label command not registered")
	}
	for _, name := range []string{"match", "label"} {
		f := c.Flags().Lookup(name)
		if f == nil {
			t.Fatalf("expected --%s flag to exist", name)
		}
	}
}
