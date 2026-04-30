package ports

import (
	"fmt"
	"sort"
	"time"
)

// Profile holds a snapshot of port usage metrics at a point in time.
type Profile struct {
	CapturedAt   time.Time
	TotalPorts   int
	ByProtocol   map[string]int
	ByState      map[string]int
	TopProcesses []ProcessPortCount
}

// ProcessPortCount pairs a process name with its open port count.
type ProcessPortCount struct {
	Process string
	Count   int
}

// Profiler builds a Profile from a slice of PortEntry values.
type Profiler struct {
	topN int
}

// NewProfiler returns a Profiler that includes the top n processes by port count.
// n must be >= 1, otherwise an error is returned.
func NewProfiler(n int) (*Profiler, error) {
	if n < 1 {
		return nil, fmt.Errorf("profiler: topN must be >= 1, got %d", n)
	}
	return &Profiler{topN: n}, nil
}

// Build analyses entries and returns a populated Profile.
func (p *Profiler) Build(entries []PortEntry) Profile {
	proto := make(map[string]int)
	state := make(map[string]int)
	proc := make(map[string]int)

	for _, e := range entries {
		proto[e.Protocol]++
		state[e.State]++
		if e.Process != "" {
			proc[e.Process]++
		}
	}

	return Profile{
		CapturedAt:   time.Now().UTC(),
		TotalPorts:   len(entries),
		ByProtocol:   proto,
		ByState:      state,
		TopProcesses: p.topProcesses(proc),
	}
}

func (p *Profiler) topProcesses(counts map[string]int) []ProcessPortCount {
	list := make([]ProcessPortCount, 0, len(counts))
	for name, c := range counts {
		list = append(list, ProcessPortCount{Process: name, Count: c})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Process < list[j].Process
	})
	if len(list) > p.topN {
		list = list[:p.topN]
	}
	return list
}
