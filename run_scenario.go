package gocuke

import (
	"reflect"
	"testing"

	messages "github.com/cucumber/messages/go/v21"
	"pgregory.net/rapid"

	"github.com/regen-network/gocuke/internal/tag"
)

func (r *docRunner) runScenario(t *testing.T, pickle *messages.Pickle, verbose bool) {
	t.Helper()

	tags := tag.NewTagsFromPickleTags(pickle.Tags)
	if r.tagExpr != nil && !tags.Match(r.tagExpr) {
		t.SkipNow()
	}

	if testing.Short() {
		if r.shortTagExpr != nil && !tags.Match(r.shortTagExpr) {
			t.SkipNow()
		}
	}

	if globalTagExpr != nil {
		if !tags.Match(globalTagExpr) {
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

	if useRapid {
		rapid.Check(t, func(rt *rapid.T) {
			t := &rapidT{rt}
			(&scenarioRunner{
				docRunner: r,
				t:         t,
				pickle:    pickle,
				stepDefs:  stepDefs,
				verbose:   verbose,
			}).runTestCase()
		})
	} else {
		(&scenarioRunner{
			docRunner: r,
			t:         t,
			pickle:    pickle,
			stepDefs:  stepDefs,
			verbose:   verbose,
		}).runTestCase()
	}
}

type scenarioRunner struct {
	*docRunner
	t        TestingT
	s        interface{}
	pickle   *messages.Pickle
	stepDefs []*stepDef
	step     *messages.PickleStep
	verbose  bool
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
		switch t := r.t.(type) {
		case *rapidT:
			defer r.runHook(hook)
		case interface{ Cleanup(func()) }:
			t.Cleanup(func() { r.runHook(hook) })
		default:
			defer r.runHook(hook)
		}
	}

	for i, step := range r.pickle.Steps {
		r.runStep(step, r.stepDefs[i])
	}
}
