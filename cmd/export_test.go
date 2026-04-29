package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestExportCommandHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"export", "--help"})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	// Help should not return an error
	_ = rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "export") {
		t.Errorf("expected 'export' in help output, got: %s", out)
	}
}

func TestExportCommandInvalidFormat(t *testing.T) {
	cmd := exportCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--format", "xml"})
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("expected error for invalid format 'xml'")
	}
	if !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("expected 'invalid format' in error, got: %v", err)
	}
}

func TestExportCmdRegistered(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "export" {
			found = true
			break
		}
	}
	if !found {
		t.Error("export command not registered on rootCmd")
	}
}

func TestExportFormatFlag(t *testing.T) {
	f := exportCmd.Flags().Lookup("format")
	if f == nil {
		t.Fatal("expected --format flag to exist")
	}
	if f.DefValue != "json" {
		t.Errorf("expected default format 'json', got '%s'", f.DefValue)
	}
}

func TestExportOutputFlag(t *testing.T) {
	f := exportCmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("expected --output flag to exist")
	}
	if f.DefValue != "" {
		t.Errorf("expected empty default output, got '%s'", f.DefValue)
	}
}

func TestExportCommandValidFormats(t *testing.T) {
	validFormats := []string{"json", "yaml"}
	for _, format := range validFormats {
		t.Run(format, func(t *testing.T) {
			cmd := exportCmd
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs([]string{"--format", format})
			err := cmd.RunE(cmd, []string{})
			if err != nil {
				t.Errorf("expected no error for valid format '%s', got: %v", format, err)
			}
		})
	}
}
