package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/portmap-cli/internal/ports"
)

var (
	classifyField    string
	classifyContains string
	classifyCategory string
	classifyJSON     bool
)

var classifyCmd = &cobra.Command{
	Use:   "classify",
	Short: "Classify active port entries into categories based on a field rule",
	Long: `Scan active ports and assign a category tag to each entry.

Example:
  portmap-cli classify --field process --contains nginx --category web`,
	RunE: runClassify,
}

func init() {
	classifyCmd.Flags().StringVar(&classifyField, "field", "process", "Field to match against (protocol, state, process)")
	classifyCmd.Flags().StringVar(&classifyContains, "contains", "", "Substring to match within the field value")
	classifyCmd.Flags().StringVar(&classifyCategory, "category", "", "Category label to assign on match")
	classifyCmd.Flags().BoolVar(&classifyJSON, "json", false, "Output results as JSON")
	_ = classifyCmd.MarkFlagRequired("contains")
	_ = classifyCmd.MarkFlagRequired("category")
	rootCmd.AddCommand(classifyCmd)
}

func runClassify(cmd *cobra.Command, _ []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	rules := []ports.ClassifierRule{
		{Field: classifyField, Contains: classifyContains, Category: classifyCategory},
	}
	classifier, err := ports.NewClassifier(rules)
	if err != nil {
		return fmt.Errorf("classifier: %w", err)
	}

	result := classifier.Classify(entries)

	if classifyJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%-8s %-6s %-20s %-20s\n", "CATEGORY", "PORT", "PROCESS", "STATE")
	for _, e := range result {
		cat := e.Tags["category"]
		fmt.Fprintf(cmd.OutOrStdout(), "%-8s %-6d %-20s %-20s\n",
			strings.ToUpper(cat), e.Port, e.Process, e.State)
	}
	return nil
}
