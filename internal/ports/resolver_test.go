package ports

import (
	"testing"
)

func TestNewResolverNoCustom(t *testing.T) {
	r := NewResolver(nil)
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestResolveWellKnownPort(t *testing.T) {
	r := NewResolver(nil)
	if got := r.Resolve(80); got != "http" {
		t.Errorf("expected http, got %s", got)
	}
	if got := r.Resolve(443); got != "https" {
		t.Errorf("expected https, got %s", got)
	}
	if got := r.Resolve(22); got != "ssh" {
		t.Errorf("expected ssh, got %s", got)
	}
}

func TestResolveUnknownPort(t *testing.T) {
	r := NewResolver(nil)
	if got := r.Resolve(19999); got != "19999" {
		t.Errorf("expected 19999, got %s", got)
	}
}

func TestResolveCustomOverridesWellKnown(t *testing.T) {
	custom := map[int]string{80: "my-app"}
	r := NewResolver(custom)
	if got := r.Resolve(80); got != "my-app" {
		t.Errorf("expected my-app, got %s", got)
	}
}

func TestResolveCustomPort(t *testing.T) {
	custom := map[int]string{9999: "my-service"}
	r := NewResolver(custom)
	if got := r.Resolve(9999); got != "my-service" {
		t.Errorf("expected my-service, got %s", got)
	}
}

func TestResolveEntry(t *testing.T) {
	r := NewResolver(nil)
	e := Entry{Port: 443}
	got := r.ResolveEntry(e)
	if got != "443 (https)" {
		t.Errorf("expected '443 (https)', got %s", got)
	}
}

func TestResolveEntryUnknown(t *testing.T) {
	r := NewResolver(nil)
	e := Entry{Port: 19999}
	got := r.ResolveEntry(e)
	if got != "19999" {
		t.Errorf("expected '19999', got %s", got)
	}
}

func TestAnnotate(t *testing.T) {
	r := NewResolver(nil)
	entries := []Entry{
		{Port: 80},
		{Port: 443},
		{Port: 19999},
		{Port: 80}, // duplicate
	}
	annotations := r.Annotate(entries)
	if annotations[80] != "http" {
		t.Errorf("expected http for port 80, got %s", annotations[80])
	}
	if annotations[443] != "https" {
		t.Errorf("expected https for port 443, got %s", annotations[443])
	}
	if _, ok := annotations[19999]; ok {
		t.Error("expected no annotation for unknown port 19999")
	}
}
