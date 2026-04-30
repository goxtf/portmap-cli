package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	compareSnapshotDir  string
	compareBaselineFile string
	compareJSONOutput   bool
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare current port state against a saved snapshot",
	Long:  "Loads a baseline snapshot and compares it to the current active ports, showing added and removed bindings.",
	RunE:  runCompare,
}

func init() {
	compareCmd.Flags().StringVarP(&compareSnapshotDir, "dir", "d", "/tmp/portmap-snapshots", "Directory containing snapshots")
	compareCmd.Flags().StringVarP(&compareBaselineFile, "baseline", "b", "", "Specific snapshot file to use as baseline (uses latest if empty)")
	compareCmd.Flags().BoolVar(&compareJSONOutput, "json", false, "Output result as JSON")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	sm, err := ports.NewSnapshotManager(compareSnapshotDir)
	if err != nil {
		return fmt.Errorf("snapshot manager: %w", err)
	}

	var baseline []ports.PortEntry
	if compareBaselineFile != "" {
		baseline, err = sm.Load(compareBaselineFile)
	} else {
		baseline, err = sm.LoadLatest()
	}
	if err != nil {
		return fmt.Errorf("loading baseline: %w", err)
	}

	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}

	current, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scanning ports: %w", err)
	}

	comparator := ports.NewComparator()
	result := comparator.Compare(baseline, current)

	if compareJSONOutput {
		return json.NewEncoder(os.Stdout).Encode(result)
	}

	fmt.Println(result.Summary())
	if len(result.Added) > 0 {
		fmt.Println("\n[+] Added:")
		for _, e := range result.Added {
			fmt.Printf("    %s/%d (%s) [%s]\n", e.Protocol, e.Port, e.Process, e.State)
		}
	}
	if len(result.Removed) > 0 {
		fmt.Println("\n[-] Removed:")
		for _, e := range result.Removed {
			fmt.Printf("    %s/%d (%s) [%s]\n", e.Protocol, e.Port, e.Process, e.State)
		}
	}
	return nil
}
