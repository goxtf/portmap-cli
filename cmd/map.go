package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var (
	mapIncludeClosed bool
	mapShowConflicts bool
	mapPort          int
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Build and display a port-to-process map",
	Long:  "Scans active connections and builds a structured port map, optionally showing conflicts or looking up a specific port.",
	RunE:  runMap,
}

func init() {
	mapCmd.Flags().BoolVar(&mapIncludeClosed, "include-closed", false, "Include CLOSE_WAIT and CLOSED entries")
	mapCmd.Flags().BoolVar(&mapShowConflicts, "conflicts", false, "Show only ports with multiple bound processes")
	mapCmd.Flags().IntVar(&mapPort, "port", 0, "Look up a specific port number")
	rootCmd.AddCommand(mapCmd)
}

func runMap(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	mapper, err := ports.NewMapper(mapIncludeClosed)
	if err != nil {
		return fmt.Errorf("mapper: %w", err)
	}
	pm := mapper.Build(entries)

	if mapPort != 0 {
		e, ok := pm.Lookup(mapPort)
		if !ok {
			fmt.Fprintf(os.Stderr, "port %d not found\n", mapPort)
			return nil
		}
		return json.NewEncoder(os.Stdout).Encode(e)
	}

	if mapShowConflicts {
		pm = pm.Conflicts()
	}

	for _, p := range pm.Ports() {
		e := pm[p]
		for _, entry := range e {
			fmt.Fprintf(os.Stdout, "%-6d %-6s %-20s %s\n", entry.Port, entry.Protocol, entry.Process, entry.State)
		}
	}
	return nil
}
