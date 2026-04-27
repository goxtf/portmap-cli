package ports

import (
	"fmt"
	"strconv"
)

// wellKnownPorts maps common port numbers to their service names.
var wellKnownPorts = map[int]string{
	20:   "ftp-data",
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	465:  "smtps",
	587:  "submission",
	993:  "imaps",
	995:  "pop3s",
	3000: "dev-server",
	3306: "mysql",
	5432: "postgresql",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	9200: "elasticsearch",
	27017: "mongodb",
}

// Resolver resolves port numbers to human-readable service names.
type Resolver struct {
	custom map[int]string
}

// NewResolver creates a new Resolver, optionally merging custom port mappings.
func NewResolver(custom map[int]string) *Resolver {
	r := &Resolver{
		custom: make(map[int]string),
	}
	for k, v := range custom {
		r.custom[k] = v
	}
	return r
}

// Resolve returns the service name for a given port number.
// Custom mappings take precedence over well-known mappings.
// If no mapping is found, the numeric port is returned as a string.
func (r *Resolver) Resolve(port int) string {
	if name, ok := r.custom[port]; ok {
		return name
	}
	if name, ok := wellKnownPorts[port]; ok {
		return name
	}
	return strconv.Itoa(port)
}

// ResolveEntry returns a display string combining port and service name.
func (r *Resolver) ResolveEntry(e Entry) string {
	service := r.Resolve(e.Port)
	if service == strconv.Itoa(e.Port) {
		return fmt.Sprintf("%d", e.Port)
	}
	return fmt.Sprintf("%d (%s)", e.Port, service)
}

// Annotate adds service name annotations to a slice of entries by
// returning a map of port -> service name for all resolved entries.
func (r *Resolver) Annotate(entries []Entry) map[int]string {
	annotations := make(map[int]string, len(entries))
	for _, e := range entries {
		if _, seen := annotations[e.Port]; !seen {
			name := r.Resolve(e.Port)
			if name != strconv.Itoa(e.Port) {
				annotations[e.Port] = name
			}
		}
	}
	return annotations
}
