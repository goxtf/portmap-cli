package ports

import (
	"fmt"
	"strings"
)

// ValidationResult holds the outcome of validating a single PortEntry.
type ValidationResult struct {
	Entry  PortEntry
	Errors []string
	Valid  bool
}

// Validator checks PortEntry values against a set of configurable rules.
type Validator struct {
	allowedProtocols map[string]bool
	allowedStates    map[string]bool
	requireProcess   bool
	maxPort          int
}

// NewValidator constructs a Validator. allowedProtocols and allowedStates are
// case-insensitive; pass nil to skip those checks. Set requireProcess to true
// to flag entries that have no associated process name.
func NewValidator(allowedProtocols, allowedStates []string, requireProcess bool, maxPort int) (*Validator, error) {
	if maxPort < 0 || maxPort > 65535 {
		return nil, fmt.Errorf("maxPort must be between 0 and 65535, got %d", maxPort)
	}
	apProto := make(map[string]bool)
	for _, p := range allowedProtocols {
		apProto[strings.ToLower(p)] = true
	}
	apState := make(map[string]bool)
	for _, s := range allowedStates {
		apState[strings.ToLower(s)] = true
	}
	return &Validator{
		allowedProtocols: apProto,
		allowedStates:    apState,
		requireProcess:   requireProcess,
		maxPort:          maxPort,
	}, nil
}

// Validate checks all entries and returns a ValidationResult per entry.
func (v *Validator) Validate(entries []PortEntry) []ValidationResult {
	results := make([]ValidationResult, 0, len(entries))
	for _, e := range entries {
		result := v.validateOne(e)
		results = append(results, result)
	}
	return results
}

func (v *Validator) validateOne(e PortEntry) ValidationResult {
	var errs []string

	if e.Port < 0 || e.Port > 65535 {
		errs = append(errs, fmt.Sprintf("port %d out of valid range 0-65535", e.Port))
	} else if v.maxPort > 0 && e.Port > v.maxPort {
		errs = append(errs, fmt.Sprintf("port %d exceeds maxPort %d", e.Port, v.maxPort))
	}

	if len(v.allowedProtocols) > 0 && !v.allowedProtocols[strings.ToLower(e.Protocol)] {
		errs = append(errs, fmt.Sprintf("protocol %q not in allowed set", e.Protocol))
	}

	if len(v.allowedStates) > 0 && !v.allowedStates[strings.ToLower(e.State)] {
		errs = append(errs, fmt.Sprintf("state %q not in allowed set", e.State))
	}

	if v.requireProcess && strings.TrimSpace(e.Process) == "" {
		errs = append(errs, "entry has no associated process")
	}

	return ValidationResult{Entry: e, Errors: errs, Valid: len(errs) == 0}
}

// Invalid returns only the results that failed validation.
func Invalid(results []ValidationResult) []ValidationResult {
	out := results[:0]
	for _, r := range results {
		if !r.Valid {
			out = append(out, r)
		}
	}
	return out
}
