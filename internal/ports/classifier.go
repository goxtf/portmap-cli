package ports

import "strings"

// ClassifierRule maps a field/value pattern to a category label.
type ClassifierRule struct {
	Field    string
	Contains string
	Category string
}

// Classifier assigns a Category field to each PortEntry based on rules.
type Classifier struct {
	rules []ClassifierRule
}

var validClassifierFields = map[string]bool{
	"protocol": true,
	"state":    true,
	"process":  true,
}

// NewClassifier creates a Classifier from the given rules.
// Returns an error if any rule references an unsupported field or has an empty category.
func NewClassifier(rules []ClassifierRule) (*Classifier, error) {
	for _, r := range rules {
		if !validClassifierFields[strings.ToLower(r.Field)] {
			return nil, &ErrInvalidField{Field: r.Field}
		}
		if strings.TrimSpace(r.Category) == "" {
			return nil, errEmptyCategory
		}
	}
	return &Classifier{rules: rules}, nil
}

// Classify applies all rules to each entry and sets the Category tag.
// The first matching rule wins. Entries with no match receive category "other".
func (c *Classifier) Classify(entries []PortEntry) []PortEntry {
	result := make([]PortEntry, len(entries))
	for i, e := range entries {
		result[i] = e
		if result[i].Tags == nil {
			result[i].Tags = make(map[string]string)
		}
		result[i].Tags["category"] = c.resolve(e)
	}
	return result
}

func (c *Classifier) resolve(e PortEntry) string {
	for _, r := range c.rules {
		var val string
		switch strings.ToLower(r.Field) {
		case "protocol":
			val = e.Protocol
		case "state":
			val = e.State
		case "process":
			val = e.Process
		}
		if strings.Contains(strings.ToLower(val), strings.ToLower(r.Contains)) {
			return r.Category
		}
	}
	return "other"
}

// errEmptyCategory is a sentinel for blank category values.
var errEmptyCategory = classifierError("classifier: category must not be empty")

type classifierError string

func (e classifierError) Error() string { return string(e) }
