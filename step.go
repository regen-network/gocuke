package gocuke

import "github.com/cucumber/common/messages/go/v19"

// Step is a special step argument type which describes the running step
// and that can be used in a step definition or hook method.
type Step interface {
	Text() string

	private()
}

type step struct {
	step *messages.PickleStep
}

func (s step) Text() string {
	return s.step.Text
}

func (s step) private() {}

var _ Step = step{}
