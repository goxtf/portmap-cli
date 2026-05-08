package ports

import (
	"testing"
)

func TestNewValidatorInvalidMaxPort(t *testing.T) {
	_, err := NewValidator(nil, nil, false, 70000)
	if err == nil {
		t.Fatal("expected error for maxPort > 65535")
	}
}

func TestNewValidatorNegativeMaxPort(t *testing.T) {
	_, err := NewValidator(nil, nil, false, -1)
	if err == nil {
		t.Fatal("expected error for negative maxPort")
	}
}

func TestNewValidatorValid(t *testing.T) {
	v, err := NewValidator([]string{"tcp", "udp"}, []string{"LISTEN"}, true, 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
}

func TestValidatePortOutOfRange(t *testing.T) {
	v, _ := NewValidator(nil, nil, false, 0)
	results := v.Validate([]PortEntry{{Port: 99999, Protocol: "tcp", State: "LISTEN"}})
	if len(results) != 1 || results[0].Valid {
		t.Fatal("expected invalid result for out-of-range port")
	}
}

func TestValidateProtocolNotAllowed(t *testing.T) {
	v, _ := NewValidator([]string{"tcp"}, nil, false, 0)
	results := v.Validate([]PortEntry{{Port: 80, Protocol: "udp", State: "LISTEN"}})
	if results[0].Valid {
		t.Fatal("expected invalid result for disallowed protocol")
	}
	if len(results[0].Errors) == 0 {
		t.Fatal("expected at least one error message")
	}
}

func TestValidateStateNotAllowed(t *testing.T) {
	v, _ := NewValidator(nil, []string{"LISTEN"}, false, 0)
	results := v.Validate([]PortEntry{{Port: 80, Protocol: "tcp", State: "CLOSE_WAIT"}})
	if results[0].Valid {
		t.Fatal("expected invalid result for disallowed state")
	}
}

func TestValidateRequireProcess(t *testing.T) {
	v, _ := NewValidator(nil, nil, true, 0)
	results := v.Validate([]PortEntry{{Port: 80, Protocol: "tcp", State: "LISTEN", Process: ""}})
	if results[0].Valid {
		t.Fatal("expected invalid when process is required but missing")
	}
}

func TestValidateMaxPort(t *testing.T) {
	v, _ := NewValidator(nil, nil, false, 1024)
	results := v.Validate([]PortEntry{{Port: 8080, Protocol: "tcp", State: "LISTEN"}})
	if results[0].Valid {
		t.Fatal("expected invalid for port exceeding maxPort")
	}
}

func TestValidateAllPassing(t *testing.T) {
	v, _ := NewValidator([]string{"tcp"}, []string{"listen"}, true, 65535)
	entries := []PortEntry{
		{Port: 80, Protocol: "tcp", State: "LISTEN", Process: "nginx"},
		{Port: 443, Protocol: "TCP", State: "listen", Process: "nginx"},
	}
	results := v.Validate(entries)
	for _, r := range results {
		if !r.Valid {
			t.Errorf("expected valid, got errors: %v", r.Errors)
		}
	}
}

func TestInvalidFiltersResults(t *testing.T) {
	v, _ := NewValidator([]string{"tcp"}, nil, false, 0)
	results := v.Validate([]PortEntry{
		{Port: 80, Protocol: "tcp", State: "LISTEN"},
		{Port: 53, Protocol: "udp", State: "LISTEN"},
	})
	invalid := Invalid(results)
	if len(invalid) != 1 {
		t.Fatalf("expected 1 invalid result, got %d", len(invalid))
	}
	if invalid[0].Entry.Protocol != "udp" {
		t.Errorf("expected udp entry to be invalid")
	}
}
