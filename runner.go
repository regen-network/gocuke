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
	initScenario         func(t TestingT) StepDefinitions
	incr                 *messages.Incrementing
	paths                []string
	parallel             bool
	stepDefs             []*stepDef
	suggestions          map[string]methodSig
	supportedSpecialArgs map[reflect.Type]specialArgGetter
	beforeHooks          []*stepDef
	afterHooks           []*stepDef
	beforeStepHooks      []*stepDef
	afterStepHooks       []*stepDef
	hooksUseRapid        bool
}

// NewRunner constructs a new Runner with the provided initScenario function.
// initScenario will be called for each test case returning a new suite instance
// for each test case which can be used for sharing state between steps. It
// is expected that the suite will retain a copy of the TestingT instance
// for usage in each step. Complex initialization should not be done in initScenario
// but rather with a Before hook.
func NewRunner(t *testing.T, initScenario func(t TestingT) StepDefinitions) *Runner {
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
		hooksUseRapid: false,
	}

	r.setupSuite(initScenario)

	return r
}

func (r *Runner) setupSuite(initScenario func(t TestingT) StepDefinitions) {
	s := initScenario(r.topLevelT)
	r.initScenario = initScenario
	r.suiteType = reflect.TypeOf(s)
	r.supportedSpecialArgs[r.suiteType] = func(runner *scenarioRunner) interface{} {
		return runner.s
	}

	addHook := func(hooks *[]*stepDef, method reflect.Method) {
		def := r.newStepDefOrHook(r.topLevelT, nil, method.Func)
		if def.usesRapid() {
			r.hooksUseRapid = true
		}
		*hooks = append(*hooks, def)
	}

	if before, ok := r.suiteType.MethodByName("Before"); ok {
		addHook(&r.beforeHooks, before)
	}

	if after, ok := r.suiteType.MethodByName("After"); ok {
		addHook(&r.afterHooks, after)
	}

	if beforeStep, ok := r.suiteType.MethodByName("BeforeStep"); ok {
		addHook(&r.beforeStepHooks, beforeStep)
	}

	if afterStep, ok := r.suiteType.MethodByName("AfterStep"); ok {
		addHook(&r.afterStepHooks, afterStep)
	}
}

// StepDefinitions is a dummy interface to mark a struct containing step definitions.
type StepDefinitions interface{}

var rapidTType = reflect.TypeOf(&rapid.T{})
