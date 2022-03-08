package gocuke

import (
	"github.com/aaronc/gocuke/internal/tag"
	"github.com/cucumber/messages-go/v16"
	"pgregory.net/rapid"
	"testing"
)

func (r *docRunner) runScenario(t *testing.T, pickle *messages.Pickle) {
	t.Helper()

	tags := tag.NewTagsFromPickleTags(pickle.Tags)
	if !r.tagExpr.Match(tags) {
		t.SkipNow()
	}

	if testing.Short() {
		if !r.shortTagExpr.Match(tags) {
			t.SkipNow()
		}
	}

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
	if cleanup, ok := s.(interface{ Cleanup() }); ok {
		if t, ok := t.(interface{ Cleanup(func()) }); ok {
			t.Cleanup(func() { cleanup.Cleanup() })
		} else {
			defer cleanup.Cleanup()
		}
	}

	for i, step := range pickle.Steps {
		r.runStep(t, step, stepDefs[i], s)
	}
}
