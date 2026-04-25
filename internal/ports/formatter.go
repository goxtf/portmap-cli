package ports

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const (
	colWidth   = 4
	minPadding = 2
)

// Formatter renders port entries to a writer in a human-readable table.
type Formatter struct {
	w io.Writer
}

// NewFormatter creates a Formatter that writes to w.
func NewFormatter(w io.Writer) *Formatter {
	return &Formatter{w: w}
}

// Render prints a formatted table of port entries.
func (f *Formatter) Render(entries []PortEntry) {
	tw := tabwriter.NewWriter(f.w, colWidth, colWidth, minPadding, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, strings.Join([]string{"PROTO", "LOCAL ADDR", "PORT", "STATE", "PID", "PROCESS"}, "\t"))
	fmt.Fprintln(tw, strings.Join([]string{"-----", "----------", "----", "-----", "---", "-------"}, "\t"))

	for _, e := range entries {
		pidStr := fmt.Sprintf("%d", e.PID)
		if e.PID == 0 {
			pidStr = "-"
		}
		process := e.Process
		if process == "" {
			process = "-"
		}
		addr := e.LocalAddr
		if addr == "" || addr == "*" {
			addr = "0.0.0.0"
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\t%s\n",
			e.Protocol,
			addr,
			e.LocalPort,
			e.State,
			pidStr,
			process,
		)
	}
}

// Summary prints a brief count summary of the entries.
func (f *Formatter) Summary(entries []PortEntry) {
	tcp, udp, other := 0, 0, 0
	for _, e := range entries {
		switch strings.ToUpper(e.Protocol) {
		case "TCP", "TCP6":
			tcp++
		case "UDP", "UDP6":
			udp++
		default:
			other++
		}
	}
	fmt.Fprintf(f.w, "\nTotal: %d  (TCP: %d, UDP: %d, Other: %d)\n",
		len(entries), tcp, udp, other)
}
