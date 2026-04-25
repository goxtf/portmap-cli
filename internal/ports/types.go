package ports

// PortEntry represents a single active port binding on the local machine.
type PortEntry struct {
	// Protocol is the network protocol (TCP or UDP).
	Protocol string

	// LocalAddress is the local IP address the port is bound to.
	LocalAddress string

	// LocalPort is the local port number as a string.
	LocalPort string

	// RemoteAddress is the remote IP address (if connected).
	RemoteAddress string

	// RemotePort is the remote port number (if connected).
	RemotePort string

	// State is the connection state (e.g., LISTEN, ESTABLISHED).
	State string

	// PID is the process identifier owning this socket.
	PID string

	// Process is the name of the process owning this socket.
	Process string
}

// String returns a compact human-readable representation of a PortEntry.
func (p PortEntry) String() string {
	addr := p.LocalAddress + ":" + p.LocalPort
	if p.RemoteAddress != "" && p.RemotePort != "" {
		addr += " -> " + p.RemoteAddress + ":" + p.RemotePort
	}
	return p.Protocol + "  " + addr + "  " + p.State + "  " + p.Process + "(" + p.PID + ")"
}
