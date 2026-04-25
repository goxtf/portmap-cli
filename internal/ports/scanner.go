package ports

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// PortEntry represents a single active port mapping with its associated process.
type PortEntry struct {
	Protocol string
	LocalAddr string
	LocalPort int
	PID       int
	Process   string
	State     string
}

// Scanner provides methods to discover active port bindings on the local machine.
type Scanner struct{}

// NewScanner creates a new Scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Scan returns all active port entries on the system.
func (s *Scanner) Scan() ([]PortEntry, error) {
	switch runtime.GOOS {
	case "linux":
		return scanLinux()
	case "darwin":
		return scanDarwin()
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func scanDarwin() ([]PortEntry, error) {
	out, err := exec.Command("lsof", "-iTCP", "-iUDP", "-n", "-P", "-sTCP:LISTEN").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsof: %w", err)
	}
	return parseLsof(string(out))
}

func scanLinux() ([]PortEntry, error) {
	out, err := exec.Command("ss", "-tulnp").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run ss: %w", err)
	}
	return parseSS(string(out))
}

func parseLsof(output string) ([]PortEntry, error) {
	var entries []PortEntry
	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		process := fields[0]
		pid, _ := strconv.Atoi(fields[1])
		proto := strings.ToUpper(fields[7])
		addr := fields[8]
		host, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			continue
		}
		port, _ := strconv.Atoi(portStr)
		entries = append(entries, PortEntry{
			Protocol:  proto,
			LocalAddr: host,
			LocalPort: port,
			PID:       pid,
			Process:   process,
			State:     "LISTEN",
		})
	}
	return entries, nil
}

func parseSS(output string) ([]PortEntry, error) {
	var entries []PortEntry
	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		proto := strings.ToUpper(fields[0])
		state := fields[1]
		localAddr := fields[4]
		host, portStr, err := net.SplitHostPort(localAddr)
		if err != nil {
			continue
		}
		port, _ := strconv.Atoi(portStr)
		pid, process := extractSSProcess(fields)
		entries = append(entries, PortEntry{
			Protocol:  proto,
			LocalAddr: host,
			LocalPort: port,
			PID:       pid,
			Process:   process,
			State:     state,
		})
	}
	return entries, nil
}

func extractSSProcess(fields []string) (int, string) {
	for _, f := range fields {
		if strings.HasPrefix(f, "users:") {
			// format: users:(("proc",pid=123,fd=4))
			parts := strings.Split(f, ",")
			var name string
			var pid int
			for _, p := range parts {
				if strings.HasPrefix(p, "((\"" ) {
					name = strings.Trim(p, `((")`)
				}
				if strings.HasPrefix(p, "pid=") {
					pid, _ = strconv.Atoi(strings.TrimPrefix(p, "pid="))
				}
			}
			return pid, name
		}
	}
	return 0, ""
}
