package ports

import (
	"fmt"
	"strings"
)

// HighlightRule defines a coloring rule for port entries.
type HighlightRule struct {
	Field   string
	Value   string
	Color   string
}

// ANSI color codes.
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

var namedColors = map[string]string{
	"red":    ColorRed,
	"green":  ColorGreen,
	"yellow": ColorYellow,
	"blue":   ColorBlue,
	"cyan":   ColorCyan,
}

// Highlighter applies color rules to port entries for terminal output.
type Highlighter struct {
	rules []HighlightRule
}

// NewHighlighter creates a Highlighter from a slice of rules.
// Returns an error if any rule references an unsupported field or unknown color.
func NewHighlighter(rules []HighlightRule) (*Highlighter, error) {
	validFields := map[string]bool{
		"protocol": true,
		"state":    true,
		"process":  true,
		"port":     true,
	}
	for _, r := range rules {
		if !validFields[strings.ToLower(r.Field)] {
			return nil, fmt.Errorf("highlighter: unsupported field %q", r.Field)
		}
		if _, ok := namedColors[strings.ToLower(r.Color)]; !ok {
			return nil, fmt.Errorf("highlighter: unknown color %q", r.Color)
		}
	}
	return &Highlighter{rules: rules}, nil
}

// Colorize returns the ANSI-colored version of a line based on the entry fields.
// The first matching rule wins.
func (h *Highlighter) Colorize(entry PortEntry, line string) string {
	for _, r := range h.rules {
		if h.matches(entry, r) {
			code := namedColors[strings.ToLower(r.Color)]
			return code + line + ColorReset
		}
	}
	return line
}

func (h *Highlighter) matches(entry PortEntry, r HighlightRule) bool {
	switch strings.ToLower(r.Field) {
	case "protocol":
		return strings.EqualFold(entry.Protocol, r.Value)
	case "state":
		return strings.EqualFold(entry.State, r.Value)
	case "process":
		return strings.EqualFold(entry.Process, r.Value)
	case "port":
		return fmt.Sprintf("%d", entry.LocalPort) == r.Value
	}
	return false
}
