package ports

// Pipeline chains multiple processing steps over a slice of PortEntry.
// Each step is a function that receives a slice and returns a (possibly
// modified) slice together with an error. Processing stops on the first
// error.
type PipelineStep func([]PortEntry) ([]PortEntry, error)

// Pipeline holds an ordered list of steps to execute.
type Pipeline struct {
	steps []PipelineStep
}

// NewPipeline constructs a Pipeline from the provided steps.
// At least one step must be supplied.
func NewPipeline(steps ...PipelineStep) (*Pipeline, error) {
	if len(steps) == 0 {
		return nil, ErrNoPipelineSteps
	}
	return &Pipeline{steps: steps}, nil
}

// Run executes every step in order, threading the output of one step into
// the input of the next. It returns the final slice or the first error
// encountered.
func (p *Pipeline) Run(entries []PortEntry) ([]PortEntry, error) {
	current := entries
	for _, step := range p.steps {
		var err error
		current, err = step(current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}

// ErrNoPipelineSteps is returned when a Pipeline is created without steps.
var ErrNoPipelineSteps = pipelineError("pipeline must have at least one step")

type pipelineError string

func (e pipelineError) Error() string { return string(e) }
