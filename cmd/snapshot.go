package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var (
	snapshotDir  string
	snapshotLoad string
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Save or diff port snapshots",
	Long:  "Save the current port state to a snapshot file, or diff two snapshots.",
	RunE:  runSnapshot,
}

func init() {
	snapshotCmd.Flags().StringVarP(&snapshotDir, "dir", "d", ".", "Directory to store snapshot files")
	snapshotCmd.Flags().StringVarP(&snapshotLoad, "diff", "c", "", "Path to a previous snapshot to diff against current state")
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("init scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan ports: %w", err)
	}

	sm, err := ports.NewSnapshotManager(snapshotDir)
	if err != nil {
		return fmt.Errorf("init snapshot manager: %w", err)
	}

	if snapshotLoad != "" {
		old, err := sm.Load(snapshotLoad)
		if err != nil {
			return fmt.Errorf("load snapshot: %w", err)
		}
		currentSnap := &ports.Snapshot{Entries: entries}
		added, removed := sm.Diff(old, currentSnap)
		fmt.Fprintf(os.Stdout, "Added (%d):\n", len(added))
		for _, e := range added {
			fmt.Fprintf(os.Stdout, "  + [%s] %s:%d (%s)\n", e.Protocol, e.LocalAddress, e.Port, e.Process)
		}
		fmt.Fprintf(os.Stdout, "Removed (%d):\n", len(removed))
		for _, e := range removed {
			fmt.Fprintf(os.Stdout, "  - [%s] %s:%d (%s)\n", e.Protocol, e.LocalAddress, e.Port, e.Process)
		}
		return nil
	}

	path, err := sm.Save(entries)
	if err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Snapshot saved: %s (%d entries)\n", path, len(entries))
	return nil
}
