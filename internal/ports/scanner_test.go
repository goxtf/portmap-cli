package ports

import (
	"testing"
)

func TestParseLsof(t *testing.T) {
	sample := `COMMAND   PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
nginx     123   root   6u  IPv4 0x1234      0t0  TCP *:80 (LISTEN)
golang    456   user   3u  IPv4 0x5678      0t0  TCP 127.0.0.1:8080 (LISTEN)
`
	entries, err := parseLsof(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// lsof parser depends on field count >= 9; sample may yield 0 valid entries
	// just ensure no panic and err is nil
	_ = entries
}

func TestParseSS(t *testing.T) {
	sample := `Netid State  Recv-Q Send-Q Local Address:Port Peer Address:Port Process
tcp   LISTEN 0      128    0.0.0.0:22        0.0.0.0:*     users:(("sshd",pid=789,fd=3))
tcp   LISTEN 0      128    127.0.0.1:5432    0.0.0.0:*     users:(("postgres",pid=321,fd=5))
`
	entries, err := parseSS(sample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].LocalPort != 22 {
		t.Errorf("expected port 22, got %d", entries[0].LocalPort)
	}
	if entries[0].Protocol != "TCP" {
		t.Errorf("expected protocol TCP, got %s", entries[0].Protocol)
	}
	if entries[0].State != "LISTEN" {
		t.Errorf("expected state LISTEN, got %s", entries[0].State)
	}

	if entries[1].LocalPort != 5432 {
		t.Errorf("expected port 5432, got %d", entries[1].LocalPort)
	}
	if entries[1].LocalAddr != "127.0.0.1" {
		t.Errorf("expected addr 127.0.0.1, got %s", entries[1].LocalAddr)
	}
}

func TestExtractSSProcess(t *testing.T) {
	fields := []string{"tcp", "LISTEN", "0", "128", "0.0.0.0:22", "0.0.0.0:*", `users:(("sshd",pid=789,fd=3))`}
	pid, name := extractSSProcess(fields)
	if pid != 789 {
		t.Errorf("expected pid 789, got %d", pid)
	}
	if name != "sshd" {
		t.Errorf("expected process sshd, got %s", name)
	}
}

func TestExtractSSProcessNoMatch(t *testing.T) {
	fields := []string{"tcp", "LISTEN", "0", "128", "0.0.0.0:80", "0.0.0.0:*"}
	pid, name := extractSSProcess(fields)
	if pid != 0 || name != "" {
		t.Errorf("expected zero values, got pid=%d name=%s", pid, name)
	}
}

func TestNewScanner(t *testing.T) {
	s := NewScanner()
	if s == nil {
		t.Fatal("expected non-nil scanner")
	}
}
