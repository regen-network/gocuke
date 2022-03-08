package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"pgregory.net/rapid"
	"reflect"
	"testing"
)

// Runner is a test runner.
type Runner struct {
	topLevelT            *testing.T
	suiteType            reflect.Type
	incr                 *messages.Incrementing
	paths                []string
	parallel             bool
	stepDefs             []*stepDef
	suggestions          map[string]methodSig
	supportedSpecialArgs map[reflect.Type]specialArgGetter
	suiteInjectors       []*suiteInjector
	beforeHooks          []*stepDef
	afterHooks           []*stepDef
	beforeStepHooks      []*stepDef
	afterStepHooks       []*stepDef
	suiteUsesRapid       bool
}

type suiteInjector struct {
	getValue specialArgGetter
	field    reflect.StructField
}

// NewRunner constructs a new Runner with the provided initScenario function.
// initScenario will be called for each test case returning a new suite instance
// for each test case which can be used for sharing state between steps. It
// is expected that the suite will retain a copy of the TestingT instance
// for usage in each step. Complex initialization should not be done in initScenario
// but rather with a Before hook.
func NewRunner(t *testing.T, stepDefinitionsType interface{}) *Runner {
	t.Helper()

	r := &Runner{
		topLevelT:   t,
		incr:        &messages.Incrementing{},
		parallel:    false,
		suggestions: map[string]methodSig{},
		supportedSpecialArgs: map[reflect.Type]specialArgGetter{
			// TestingT
			reflect.TypeOf((*TestingT)(nil)).Elem(): func(runner *scenarioRunner) interface{} {
				return runner.t
			},
			// *rapid.T
			reflect.TypeOf(&rapid.T{}): func(runner *scenarioRunner) interface{} {
				if t, ok := runner.t.(*rapid.T); ok {
					return t
				}
				runner.t.Fatalf("expected %T, but got %T", &rapid.T{}, runner.t)
				return nil
			},
			// Scenario
			reflect.TypeOf((*Scenario)(nil)).Elem(): func(runner *scenarioRunner) interface{} {
				return scenario{runner.pickle}
			},
		},
		suiteUsesRapid: false,
	}

	r.registerSuite(stepDefinitionsType)

	return r
}

func (r *Runner) registerSuite(stepDefinitionsType interface{}) *Runner {
	r.topLevelT.Helper()

	typ := reflect.TypeOf(stepDefinitionsType)
	r.suiteType = typ
	kind := typ.Kind()

	suiteElemType := r.suiteType
	if kind == reflect.Ptr {
		suiteElemType = suiteElemType.Elem()
	}

	if suiteElemType.Kind() != reflect.Struct {
		r.topLevelT.Fatalf("expected a struct or a pointer to a struct, got %T", stepDefinitionsType)
	}

	for i := 0; i < suiteElemType.NumField(); i++ {
		field := suiteElemType.Field(i)
		if !field.IsExported() {
			continue
		}

		if getter, ok := r.supportedSpecialArgs[field.Type]; ok {
			r.suiteInjectors = append(r.suiteInjectors, &suiteInjector{getValue: getter, field: field})
		}
	}

	r.supportedSpecialArgs[r.suiteType] = func(runner *scenarioRunner) interface{} {
		return runner.s
	}

	if before, ok := r.suiteType.MethodByName("Before"); ok {
		r.addHook(&r.beforeHooks, before.Func)
	}

	if after, ok := r.suiteType.MethodByName("After"); ok {
		r.addHook(&r.afterHooks, after.Func)
	}

	if beforeStep, ok := r.suiteType.MethodByName("BeforeStep"); ok {
		r.addHook(&r.beforeStepHooks, beforeStep.Func)
	}

	if afterStep, ok := r.suiteType.MethodByName("AfterStep"); ok {
		r.addHook(&r.afterStepHooks, afterStep.Func)
	}

	return r
}

var rapidTType = reflect.TypeOf(&rapid.T{})
