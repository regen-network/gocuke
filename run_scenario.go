package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"pgregory.net/rapid"
	"testing"
)

func (r *docRunner) runScenario(t *testing.T, pickle *messages.Pickle) {
	t.Helper()
	stepDefs := make([]*stepDef, len(pickle.Steps))
	for i, step := range pickle.Steps {
		stepDefs[i] = r.findStep(t, step)
	}

	if t.Failed() {
		return
	}

	haveRapid := false
	for _, def := range stepDefs {
		for _, arg := range def.specialArgs {
			if arg.typ == rapidTType {
				haveRapid = true
				break
			}
		}
		if haveRapid {
			break
		}
	}

	t.Run(pickle.Name, func(t *testing.T) {
		if r.parallel {
			t.Parallel()
		}

		if haveRapid {
			rapid.Check(t, func(t *rapid.T) {
				(&scenarioRunner{
					docRunner: r,
					t:         t,
					pickle:    pickle,
					stepDefs:  stepDefs,
				}).runTestCase()
			})
		} else {
			(&scenarioRunner{
				docRunner: r,
				t:         t,
				pickle:    pickle,
				stepDefs:  stepDefs,
			}).runTestCase()
		}
	})
}

type scenarioRunner struct {
	*docRunner
	t        TestingT
	s        Suite
	pickle   *messages.Pickle
	stepDefs []*stepDef
}

func (r *scenarioRunner) runTestCase() {
	r.s = r.initScenario(r.t)
	for _, hook := range r.beforeHooks {
		r.runHook(hook)
	}

	for _, hook := range r.afterHooks {
		if t, ok := r.t.(interface{ Cleanup(func()) }); ok {
			t.Cleanup(func() { r.runHook(hook) })
		} else {
			defer r.runHook(hook)
		}
	}

	for i, step := range r.pickle.Steps {
		r.runStep(step, r.stepDefs[i])
	}
}
