package gocuke

import (
	"reflect"
	"sync"
	"testing"

	messages "github.com/cucumber/messages/go/v22"
	tag "github.com/cucumber/tag-expressions/go/v6"
	"pgregory.net/rapid"
)

// Runner is a test runner.
type Runner struct {
	topLevelT            *testing.T
	suiteType            reflect.Type
	incr                 *messages.Incrementing
	paths                []string
	parallel             bool
	stepDefs             []*stepDef
	stepDefsMutex 	  	 *sync.RWMutex
	haveSuggestion       map[string]bool
	suggestions          []methodSig
	supportedSpecialArgs map[reflect.Type]specialArgGetter
	suiteInjectors       []*suiteInjector
	beforeHooks          []*stepDef
	afterHooks           []*stepDef
	beforeStepHooks      []*stepDef
	afterStepHooks       []*stepDef
	suiteUsesRapid       bool
	tagExpr              tag.Evaluatable
	shortTagExpr         tag.Evaluatable
}

type suiteInjector struct {
	getValue specialArgGetter
	field    reflect.StructField
}

// NewRunner constructs a new Runner with the provided suite type instance.
// Suite type is expected to be a pointer to a struct or a struct.
// A new instance of suiteType will be constructed for every scenario.
//
// The following special argument types will be injected into exported fields of
// the suite type struct: TestingT, Scenario, Step, *rapid.T.
//
// Methods defined on the suite type will be auto-registered as step definitions
// if they correspond to the expected method name for a step. Method
// parameters can start with the special argument types listed above and must
// be followed by step argument types for each captured step argument and
// DocString or DataTable at the end if the step uses one of these.
// Valid step argument types are int64, string, *big.Int and *apd.Decimal.
//
// The methods Before, After, BeforeStep and AfterStep will be recognized
// as hooks and can take the special argument types listed above.
func NewRunner(t *testing.T, suiteType interface{}) *Runner {
	t.Helper()

	initGlobalTagExpr()

	r := &Runner{
		topLevelT:      t,
		incr:           &messages.Incrementing{},
		parallel:       true,
		haveSuggestion: map[string]bool{},
		supportedSpecialArgs: map[reflect.Type]specialArgGetter{
			// TestingT
			reflect.TypeOf((*TestingT)(nil)).Elem(): func(runner *scenarioRunner) interface{} {
				return runner.t
			},
			// *rapid.T
			reflect.TypeOf(&rapid.T{}): func(runner *scenarioRunner) interface{} {
				if t, ok := runner.t.(*rapidT); ok {
					return t.T
				}
				runner.t.Fatalf("expected %T, but got %T", &rapid.T{}, runner.t)
				return nil
			},
			// Scenario
			reflect.TypeOf((*Scenario)(nil)).Elem(): func(runner *scenarioRunner) interface{} {
				return scenario{runner.pickle}
			},
			// Step
			reflect.TypeOf((*Step)(nil)).Elem(): func(runner *scenarioRunner) interface{} {
				return step{runner.step}
			},
		},
		suiteUsesRapid: false,
		stepDefsMutex: &sync.RWMutex{},
	}

	r.registerSuite(suiteType)

	return r
}

func (r *Runner) registerSuite(suiteType interface{}) *Runner {
	r.topLevelT.Helper()

	typ := reflect.TypeOf(suiteType)
	r.suiteType = typ
	kind := typ.Kind()

	suiteElemType := r.suiteType
	if kind == reflect.Ptr {
		suiteElemType = suiteElemType.Elem()
	}

	if suiteElemType.Kind() != reflect.Struct {
		r.topLevelT.Fatalf("expected a struct or a pointer to a struct, got %T", suiteType)
	}

	for i := 0; i < suiteElemType.NumField(); i++ {
		field := suiteElemType.Field(i)
		if !field.IsExported() {
			continue
		}

		for specialArgType, getter := range r.supportedSpecialArgs {
			if field.Type.AssignableTo(specialArgType) {
				r.suiteInjectors = append(r.suiteInjectors, &suiteInjector{getValue: getter, field: field})
				if field.Type == rapidTType {
					r.suiteUsesRapid = true
				}
				break
			}
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
