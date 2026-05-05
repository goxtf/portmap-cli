package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	scoreWeightPriv    float64
	scoreWeightProcess float64
	scoreWeightTCP     float64
	scoreWeightListen  float64
	scoreJSONOutput    bool
	scoreTopN          int
)

func init() {
	scoreCmd := &cobra.Command{
		Use:   "score",
		Short: "Rank active port entries by weighted relevance score",
		Long:  "Scores each port entry across multiple dimensions (privileged port, process presence, protocol, state) and outputs them ranked highest-first.",
		RunE:  runScore,
	}
	scoreCmd.Flags().Float64Var(&scoreWeightPriv, "weight-priv", 3.0, "Weight for ports below 1024")
	scoreCmd.Flags().Float64Var(&scoreWeightProcess, "weight-process", 2.0, "Weight for entries with a bound process")
	scoreCmd.Flags().Float64Var(&scoreWeightTCP, "weight-tcp", 1.0, "Weight for TCP protocol entries")
	scoreCmd.Flags().Float64Var(&scoreWeightListen, "weight-listen", 1.5, "Weight for entries in LISTEN state")
	scoreCmd.Flags().BoolVar(&scoreJSONOutput, "json", false, "Output results as JSON")
	scoreCmd.Flags().IntVar(&scoreTopN, "top", 0, "Limit output to top N results (0 = all)")
	rootCmd.AddCommand(scoreCmd)
}

func runScore(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("score: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("score: %w", err)
	}

	weights := ports.ScoreWeights{
		PortBelow1024: scoreWeightPriv,
		HasProcess:    scoreWeightProcess,
		TCPProtocol:   scoreWeightTCP,
		ListenState:   scoreWeightListen,
	}
	scorer, err := ports.NewScorer(weights)
	if err != nil {
		return fmt.Errorf("score: %w", err)
	}

	ranked := scorer.Rank(entries)
	if scoreTopN > 0 && scoreTopN < len(ranked) {
		ranked = ranked[:scoreTopN]
	}

	if scoreJSONOutput {
		return json.NewEncoder(os.Stdout).Encode(ranked)
	}
	fmt.Fprintf(os.Stdout, "%-6s %-8s %-14s %-12s %s\n", "SCORE", "PORT", "PROTOCOL", "STATE", "PROCESS")
	for _, se := range ranked {
		fmt.Fprintf(os.Stdout, "%-6.1f %-8d %-14s %-12s %s\n",
			se.Score, se.Entry.Port, se.Entry.Protocol, se.Entry.State, se.Entry.Process)
	}
	return nil
}
