package ports

import (
	"strings"
	"testing"
)

func TestNewHighlighterValidRules(t *testing.T) {
	rules := []HighlightRule{
		{Field: "protocol", Value: "TCP", Color: "green"},
		{Field: "state", Value: "LISTEN", Color: "cyan"},
	}
	h, err := NewHighlighter(rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Highlighter")
	}
}

func TestNewHighlighterInvalidField(t *testing.T) {
	rules := []HighlightRule{
		{Field: "unknown", Value: "foo", Color: "red"},
	}
	_, err := NewHighlighter(rules)
	if err == nil {
		t.Fatal("expected error for invalid field")
	}
	if !strings.Contains(err.Error(), "unsupported field") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewHighlighterInvalidColor(t *testing.T) {
	rules := []HighlightRule{
		{Field: "state", Value: "LISTEN", Color: "purple"},
	}
	_, err := NewHighlighter(rules)
	if err == nil {
		t.Fatal("expected error for unknown color")
	}
	if !strings.Contains(err.Error(), "unknown color") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestColorizeMatchesProtocol(t *testing.T) {
	h, _ := NewHighlighter([]HighlightRule{
		{Field: "protocol", Value: "TCP", Color: "green"},
	})
	entry := PortEntry{Protocol: "TCP", LocalPort: 80, State: "LISTEN", Process: "nginx"}
	result := h.Colorize(entry, "nginx 80 TCP")
	if !strings.Contains(result, ColorGreen) {
		t.Errorf("expected green color code in output, got: %q", result)
	}
	if !strings.Contains(result, ColorReset) {
		t.Errorf("expected reset code in output")
	}
}

func TestColorizeNoMatch(t *testing.T) {
	h, _ := NewHighlighter([]HighlightRule{
		{Field: "protocol", Value: "UDP", Color: "red"},
	})
	entry := PortEntry{Protocol: "TCP", LocalPort: 443, State: "LISTEN", Process: "nginx"}
	line := "nginx 443 TCP"
	result := h.Colorize(entry, line)
	if result != line {
		t.Errorf("expected unchanged line, got: %q", result)
	}
}

func TestColorizeFirstRuleWins(t *testing.T) {
	h, _ := NewHighlighter([]HighlightRule{
		{Field: "state", Value: "LISTEN", Color: "blue"},
		{Field: "state", Value: "LISTEN", Color: "red"},
	})
	entry := PortEntry{Protocol: "TCP", LocalPort: 8080, State: "LISTEN", Process: "app"}
	result := h.Colorize(entry, "app 8080")
	if !strings.Contains(result, ColorBlue) {
		t.Errorf("expected blue (first rule), got: %q", result)
	}
	if strings.Contains(result, ColorRed) {
		t.Errorf("did not expect red (second rule), got: %q", result)
	}
}

func TestColorizeByPort(t *testing.T) {
	h, _ := NewHighlighter([]HighlightRule{
		{Field: "port", Value: "22", Color: "yellow"},
	})
	entry := PortEntry{Protocol: "TCP", LocalPort: 22, State: "LISTEN", Process: "sshd"}
	result := h.Colorize(entry, "sshd 22")
	if !strings.Contains(result, ColorYellow) {
		t.Errorf("expected yellow for port 22, got: %q", result)
	}
}
