package ports

import "sort"

// Stats holds aggregated statistics about port mappings.
type Stats struct {
	TotalEntries    int
	ByProtocol      map[string]int
	ByState         map[string]int
	TopProcesses    []ProcessCount
	UniqueProcesses int
	UniquePorts     int
}

// ProcessCount pairs a process name with its port count.
type ProcessCount struct {
	Process string
	Count   int
}

// StatsCollector computes statistics from port entries.
type StatsCollector struct{}

// NewStatsCollector returns a new StatsCollector.
func NewStatsCollector() *StatsCollector {
	return &StatsCollector{}
}

// Collect computes statistics from the provided entries.
func (s *StatsCollector) Collect(entries []PortEntry) Stats {
	stats := Stats{
		TotalEntries: len(entries),
		ByProtocol:   make(map[string]int),
		ByState:      make(map[string]int),
	}

	processCounts := make(map[string]int)
	portSet := make(map[int]struct{})

	for _, e := range entries {
		stats.ByProtocol[e.Protocol]++
		stats.ByState[e.State]++
		processCounts[e.Process]++
		portSet[e.Port] = struct{}{}
	}

	stats.UniqueProcesses = len(processCounts)
	stats.UniquePorts = len(portSet)
	stats.TopProcesses = topN(processCounts, 5)

	return stats
}

func topN(counts map[string]int, n int) []ProcessCount {
	list := make([]ProcessCount, 0, len(counts))
	for proc, cnt := range counts {
		list = append(list, ProcessCount{Process: proc, Count: cnt})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Count != list[j].Count {
			return list[i].Count > list[j].Count
		}
		return list[i].Process < list[j].Process
	})
	if len(list) > n {
		return list[:n]
	}
	return list
}
