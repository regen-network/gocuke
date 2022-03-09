package gocuke

import "github.com/cucumber/messages-go/v16"

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
