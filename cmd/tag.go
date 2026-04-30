package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"portmap-cli/internal/ports"
)

var (
	tagField  string
	tagMatch  string
	tagKey    string
	tagValue  string
	tagJSON   bool
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag active port entries based on field matching rules",
	Long:  `Scan active ports and apply a label (tag) to entries that match the given field and value.`,
	RunE:  runTag,
}

func init() {
	tagCmd.Flags().StringVar(&tagField, "field", "protocol", "Field to match on (protocol, state, process, port)")
	tagCmd.Flags().StringVar(&tagMatch, "match", "", "Value to match against the field")
	tagCmd.Flags().StringVar(&tagKey, "key", "", "Tag key to apply")
	tagCmd.Flags().StringVar(&tagValue, "value", "", "Tag value to apply")
	tagCmd.Flags().BoolVar(&tagJSON, "json", false, "Output result as JSON")
	_ = tagCmd.MarkFlagRequired("match")
	_ = tagCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(tagCmd)
}

func runTag(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner init: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	rules := []ports.TagRule{
		{Field: tagField, Match: tagMatch, Key: tagKey, Value: tagValue},
	}
	tagger, err := ports.NewTagger(rules)
	if err != nil {
		return fmt.Errorf("tagger init: %w", err)
	}

	tagged := tagger.Tag(entries)

	if tagJSON {
		type result struct {
			Entry ports.PortEntry `json:"entry"`
			Tags  []ports.Tag     `json:"tags"`
		}
		var out []result
		for i, e := range entries {
			if tags, ok := tagged[i]; ok {
				out = append(out, result{Entry: e, Tags: tags})
			}
		}
		return json.NewEncoder(os.Stdout).Encode(out)
	}

	for i, tags := range tagged {
		e := entries[i]
		for _, tag := range tags {
			fmt.Printf("%s\t%d\t%s\t%s\t[%s=%s]\n",
				e.Protocol, e.Port, e.State, e.Process, tag.Key, tag.Value)
		}
	}
	return nil
}
