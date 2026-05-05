package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	labelField  string
	labelMatch  string
	labelValue  string
	labelJSON   bool
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Attach labels to port entries matching a rule",
	Long:  `Scan active ports and attach a text label to entries where a specified field matches a given value.`,
	RunE:  runLabel,
}

func init() {
	labelCmd.Flags().StringVarP(&labelField, "field", "f", "protocol", "Field to match (protocol, state, process, port)")
	labelCmd.Flags().StringVarP(&labelMatch, "match", "m", "", "Value to match against the field")
	labelCmd.Flags().StringVarP(&labelValue, "label", "l", "", "Label to attach to matching entries")
	labelCmd.Flags().BoolVar(&labelJSON, "json", false, "Output results as JSON")
	_ = labelCmd.MarkFlagRequired("match")
	_ = labelCmd.MarkFlagRequired("label")
	rootCmd.AddCommand(labelCmd)
}

func runLabel(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	rules := []ports.LabelRule{
		{Field: labelField, Match: labelMatch, Label: labelValue},
	}
	lb, err := ports.NewLabeler(rules)
	if err != nil {
		return fmt.Errorf("labeler: %w", err)
	}
	labeled := lb.Apply(entries)

	if labelJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(labeled)
	}

	for _, e := range labeled {
		labelList := []string{}
		for k, v := range e.Tags {
			if strings.HasPrefix(k, "label:") {
				labelList = append(labelList, v)
			}
		}
		tagStr := ""
		if len(labelList) > 0 {
			tagStr = " [" + strings.Join(labelList, ",") + "]"
		}
		fmt.Fprintf(os.Stdout, "%-6s %5d  %-12s  %-10s%s\n",
			e.Protocol, e.Port, e.Process, e.State, tagStr)
	}
	return nil
}
