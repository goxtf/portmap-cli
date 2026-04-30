package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var profileTopN int
var profileJSON bool

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Display a statistical profile of active port mappings",
	RunE:  runProfile,
}

func init() {
	profileCmd.Flags().IntVarP(&profileTopN, "top", "n", 5, "Number of top processes to display")
	profileCmd.Flags().BoolVar(&profileJSON, "json", false, "Output profile as JSON")
	rootCmd.AddCommand(profileCmd)
}

func runProfile(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("profile: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("profile: scan failed: %w", err)
	}

	profiler, err := ports.NewProfiler(profileTopN)
	if err != nil {
		return fmt.Errorf("profile: %w", err)
	}

	profile := profiler.Build(entries)

	if profileJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(profile)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Captured At:\t%s\n", profile.CapturedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Total Ports:\t%d\n", profile.TotalPorts)
	fmt.Fprintln(w, "\nBy Protocol:")
	for k, v := range profile.ByProtocol {
		fmt.Fprintf(w, "  %s\t%d\n", k, v)
	}
	fmt.Fprintln(w, "\nBy State:")
	for k, v := range profile.ByState {
		fmt.Fprintf(w, "  %s\t%d\n", k, v)
	}
	fmt.Fprintf(w, "\nTop %d Processes:\n", profileTopN)
	for _, p := range profile.TopProcesses {
		fmt.Fprintf(w, "  %s\t%d\n", p.Process, p.Count)
	}
	return w.Flush()
}
