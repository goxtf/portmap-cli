package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPipelineCommandHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"pipeline", "--help"})
	_ = rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "pipeline") {
		t.Errorf("help output missing 'pipeline': %s", out)
	}
}

func TestPipelineCmdRegistered(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Name() == "pipeline" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("pipeline command not registered on root")
	}
}

// findPipelineCmd is a helper that locates the pipeline subcommand from the
// root command's children, returning nil if it is not registered.
func findPipelineCmd() *cobra.Command {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "pipeline" {
			return c
		}
	}
	return nil
}

func TestPipelineProtocolFlag(t *testing.T) {
	cmd := findPipelineCmd()
	if cmd == nil {
		t.Fatal("pipeline command not found")
	}
	f := cmd.Flags().Lookup("protocol")
	if f == nil {
		t.Fatal("--protocol flag not registered")
	}
}

func TestPipelineDedupFlag(t *testing.T) {
	cmd := findPipelineCmd()
	if cmd == nil {
		t.Fatal("pipeline command not found")
	}
	f := cmd.Flags().Lookup("dedup")
	if f == nil {
		t.Fatal("--dedup flag not registered")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}

func TestPipelineSortFlag(t *testing.T) {
	cmd := findPipelineCmd()
	if cmd == nil {
		t.Fatal("pipeline command not found")
	}
	f := cmd.Flags().Lookup("sort")
	if f == nil {
		t.Fatal("--sort flag not registered")
	}
	if f.DefValue != "port" {
		t.Errorf("expected default 'port', got %s", f.DefValue)
	}
}
