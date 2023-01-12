package gocuke

import (
	"reflect"
	"testing"
	"time"

	messages "github.com/cucumber/messages/go/v21"
	"pgregory.net/rapid"

	"github.com/regen-network/gocuke/internal/tag"
)

func (r *docRunner) runScenario(t *testing.T, pickle *messages.Pickle) {
	t.Helper()

	if r.reporter != nil {
		r.reporter.Report(&messages.Envelope{Pickle: pickle})
	}

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
	missingSteps := false
	for i, step := range pickle.Steps {
		step := r.findStep(t, step)
		if step == nil {
			missingSteps = true
		}
		stepDefs[i] = step
	}

	if missingSteps {
		t.FailNow()
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

	var testCaseId string
	if r.reporter != nil {
		testCaseId = newId()
		r.reporter.Report(&messages.Envelope{TestCase: &messages.TestCase{
			Id:        testCaseId,
			PickleId:  pickle.Id,
			TestSteps: nil,
		}})
	}

	if useRapid {
		var attempt int64 = 0
		rapid.Check(t, func(t *rapid.T) {
			(&scenarioRunner{
				docRunner:  r,
				t:          t,
				pickle:     pickle,
				stepDefs:   stepDefs,
				testCaseId: testCaseId,
				attempt:    attempt,
			}).runTestCase()
			attempt++
		})
	} else {
		(&scenarioRunner{
			docRunner:  r,
			t:          t,
			pickle:     pickle,
			testCaseId: testCaseId,
			stepDefs:   stepDefs,
		}).runTestCase()
	}
}

type scenarioRunner struct {
	*docRunner
	t                 TestingT
	s                 interface{}
	pickle            *messages.Pickle
	stepDefs          []*stepDef
	step              *messages.PickleStep
	attempt           int64
	testCaseId        string
	testCaseStartedId string
	testStepId        string
}

func (r *scenarioRunner) runTestCase() {
	r.t.Helper()

	if r.reporter != nil {
		timestamp := messages.GoTimeToTimestamp(time.Now())
		r.testCaseStartedId = newId()
		r.reporter.Report(&messages.Envelope{TestCaseStarted: &messages.TestCaseStarted{
			Attempt:    r.attempt,
			Id:         r.testCaseStartedId,
			TestCaseId: r.testCaseId,
			Timestamp:  &timestamp,
		}})
		if _, ok := r.t.(*testing.T); ok {
			r.t = &testingTWrapper{TestingT: r.t}
		}
		defer func() {
			timestamp := messages.GoTimeToTimestamp(time.Now())
			r.reporter.Report(&messages.Envelope{TestCaseFinished: &messages.TestCaseFinished{
				TestCaseStartedId: r.testCaseStartedId,
				Timestamp:         &timestamp,
			}})
		}()
	}

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
