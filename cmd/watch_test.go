package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestWatchCommandHelp(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"watch", "--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "watch") {
		t.Errorf("expected 'watch' in help output, got: %s", output)
	}
	if !strings.Contains(output, "interval") {
		t.Errorf("expected 'interval' flag in help output, got: %s", output)
	}
}

func TestWatchCmdRegistered(t *testing.T) {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "watch" {
			return
		}
	}
	t.Error("watch command not registered on root")
}

func TestWatchIntervalFlag(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"watch"})
	if err != nil || cmd == nil {
		t.Fatal("watch command not found")
	}
	flag := cmd.Flags().Lookup("interval")
	if flag == nil {
		t.Fatal("interval flag not found on watch command")
	}
	if flag.DefValue != "2" {
		t.Errorf("expected default interval 2, got %s", flag.DefValue)
	}
}

func TestWatchIntervalFlagShorthand(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"watch"})
	if err != nil || cmd == nil {
		t.Fatal("watch command not found")
	}
	flag := cmd.Flags().ShorthandLookup("i")
	if flag == nil {
		t.Error("expected shorthand -i for interval flag")
	}
}
