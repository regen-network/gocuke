package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"pgregory.net/rapid"
	"testing"
)

func (r *Runner) runScenario(t *testing.T, pickle *messages.Pickle) {
	stepDefs := make([]*stepDef, len(pickle.Steps))
	for i, step := range pickle.Steps {
		stepDefs[i] = r.findStep(t, step)
	}

	if t.Failed() {
		return
	}

	haveRapid := false
	for _, def := range stepDefs {
		if def.hasRapid {
			haveRapid = true
			break
		}
	}

	t.Run(pickle.Name, func(t *testing.T) {
		if r.parallel {
			t.Parallel()
		}

		if haveRapid {
			rapid.Check(t, func(t *rapid.T) {
				r.runSteps(t, pickle, stepDefs)
			})
		} else {
			r.runSteps(t, pickle, stepDefs)
		}
	})
}

func (r *Runner) runSteps(t TestingT, pickle *messages.Pickle, stepDefs []*stepDef) {
	s := r.initScenario(t)
	for i, step := range pickle.Steps {
		r.runStep(t, step, stepDefs[i], s)
	}
}
