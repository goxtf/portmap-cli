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
	validateProtocols string
	validateStates    string
	validateRequireProc bool
	validateMaxPort  int
	validateJSONOut  bool
)

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate active port entries against configurable rules",
		Long:  "Scan active port mappings and report entries that violate the specified validation rules.",
		RunE:  runValidate,
	}
	validateCmd.Flags().StringVar(&validateProtocols, "protocols", "", "comma-separated list of allowed protocols (e.g. tcp,udp)")
	validateCmd.Flags().StringVar(&validateStates, "states", "", "comma-separated list of allowed states (e.g. LISTEN)")
	validateCmd.Flags().BoolVar(&validateRequireProc, "require-process", false, "flag entries with no associated process")
	validateCmd.Flags().IntVar(&validateMaxPort, "max-port", 65535, "maximum allowed port number")
	validateCmd.Flags().BoolVar(&validateJSONOut, "json", false, "output results as JSON")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	scanner, err := ports.NewScanner()
	if err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	entries, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	var protos, states []string
	if validateProtocols != "" {
		protos = splitCSV(validateProtocols)
	}
	if validateStates != "" {
		states = splitCSV(validateStates)
	}

	validator, err := ports.NewValidator(protos, states, validateRequireProc, validateMaxPort)
	if err != nil {
		return fmt.Errorf("validator: %w", err)
	}

	results := validator.Validate(entries)
	invalid := ports.Invalid(results)

	if validateJSONOut {
		return json.NewEncoder(os.Stdout).Encode(invalid)
	}

	if len(invalid) == 0 {
		fmt.Println("All entries passed validation.")
		return nil
	}
	fmt.Fprintf(os.Stdout, "%d invalid entries:\n", len(invalid))
	for _, r := range invalid {
		e := r.Entry
		fmt.Fprintf(os.Stdout, "  [%s/%d %s %s] %s\n",
			e.Protocol, e.Port, e.State, e.Process,
			strings.Join(r.Errors, "; "))
	}
	return nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
