package ports

import "fmt"

// InspectionResult holds detailed information about a single port entry.
type InspectionResult struct {
	Entry       PortEntry
	ServiceName string
	Tags        []string
	Score       float64
	Conflicts   []PortEntry
	Note        string
}

// Inspector provides deep inspection of a single port entry in context.
type Inspector struct {
	resolver *Resolver
	scorer   *Scorer
	tagger   *Tagger
	mapper   *Mapper
}

// NewInspector creates an Inspector. resolver, scorer, tagger, and mapper may be nil.
func NewInspector(resolver *Resolver, scorer *Scorer, tagger *Tagger, mapper *Mapper) (*Inspector, error) {
	return &Inspector{
		resolver: resolver,
		scorer:   scorer,
		tagger:   tagger,
		mapper:   mapper,
	}, nil
}

// Inspect returns an InspectionResult for the given entry within the provided context slice.
func (i *Inspector) Inspect(entry PortEntry, all []PortEntry) InspectionResult {
	result := InspectionResult{Entry: entry}

	if i.resolver != nil {
		result.ServiceName = i.resolver.Resolve(entry.Port)
	}

	if i.scorer != nil {
		result.Score = i.scorer.Score(entry)
	}

	if i.tagger != nil {
		tagged := i.tagger.Apply([]PortEntry{entry})
		if len(tagged) > 0 {
			result.Tags = tagged[0].Tags
		}
	}

	if i.mapper != nil {
		result.Conflicts = i.mapper.Conflicts(entry.Port)
	}

	result.Note = i.buildNote(entry, all)
	return result
}

func (i *Inspector) buildNote(entry PortEntry, all []PortEntry) string {
	count := 0
	for _, e := range all {
		if e.Port == entry.Port {
			count++
		}
	}
	if count > 1 {
		return fmt.Sprintf("port %d is shared by %d entries", entry.Port, count)
	}
	if entry.Process == "" {
		return "no owning process detected"
	}
	return ""
}
