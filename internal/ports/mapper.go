package ports

import "sort"

// PortMap holds a mapping from port number to a list of entries bound to it.
type PortMap map[int][]PortEntry

// Mapper builds a structured map of port → entries for quick lookup.
type Mapper struct {
	includeClosed bool
}

// NewMapper creates a Mapper. When includeClosed is false, entries whose
// State is "CLOSE_WAIT" or "CLOSED" are omitted from the map.
func NewMapper(includeClosed bool) (*Mapper, error) {
	return &Mapper{includeClosed: includeClosed}, nil
}

// Build converts a flat slice of PortEntry values into a PortMap.
func (m *Mapper) Build(entries []PortEntry) PortMap {
	pm := make(PortMap)
	for _, e := range entries {
		if !m.includeClosed {
			if e.State == "CLOSE_WAIT" || e.State == "CLOSED" {
				continue
			}
		}
		pm[e.Port] = append(pm[e.Port], e)
	}
	return pm
}

// Ports returns all port numbers present in the map, sorted ascending.
func (pm PortMap) Ports() []int {
	keys := make([]int, 0, len(pm))
	for k := range pm {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// Lookup returns the entries for a specific port, and whether the port exists.
func (pm PortMap) Lookup(port int) ([]PortEntry, bool) {
	e, ok := pm[port]
	return e, ok
}

// Conflicts returns ports that have more than one entry bound to them.
func (pm PortMap) Conflicts() PortMap {
	out := make(PortMap)
	for port, entries := range pm {
		if len(entries) > 1 {
			out[port] = entries
		}
	}
	return out
}
