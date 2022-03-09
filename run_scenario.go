package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"github.com/regen-network/gocuke/internal/tag"
	"pgregory.net/rapid"
	"reflect"
	"testing"
)

func (r *docRunner) runScenario(t *testing.T, pickle *messages.Pickle) {
	t.Helper()

	tags := tag.NewTagsFromPickleTags(pickle.Tags)
	if r.tagExpr != nil && !r.tagExpr.Match(tags) {
		t.SkipNow()
	}

	if testing.Short() {
		if r.shortTagExpr != nil && !r.shortTagExpr.Match(tags) {
			t.SkipNow()
		}
	}

	if globalTagExpr != nil {
		if !globalTagExpr.Match(tags) {
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

	useRapid := r.suiteUsesRapid
	if !useRapid {
		for _, def := range stepDefs {
			if def.usesRapid() {
				useRapid = true
				break
			}
		}
	}

	t.Run(pickle.Name, func(t *testing.T) {
		t.Helper()

		if r.parallel {
			t.Parallel()
		}

		if useRapid {
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
	s        interface{}
	pickle   *messages.Pickle
	stepDefs []*stepDef
	step     *messages.PickleStep
}

func (r *scenarioRunner) runTestCase() {
	r.t.Helper()

	var val reflect.Value
	needPtr := r.suiteType.Kind() == reflect.Ptr
	if needPtr {
		val = reflect.New(r.suiteType.Elem())
	} else {
		val = reflect.New(r.suiteType)
	}

	for _, injector := range r.suiteInjectors {
		val.Elem().FieldByName(injector.field.Name).Set(reflect.ValueOf(injector.getValue(r)))
	}

	if needPtr {
		r.s = val.Interface()
	} else {
		r.s = val.Elem().Interface()
	}

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
