package gocuke

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"testing"

	"gotest.tools/v3/assert"
)

type stepDef struct {
	regex       *regexp.Regexp
	theFunc     reflect.Value
	specialArgs []*specialArg
	funcLoc     string
}

type specialArg struct {
	typ      reflect.Type
	getValue specialArgGetter
}

type specialArgGetter func(*scenarioRunner) interface{}

// Step can be used to manually register a step with the runner. step should
// be a string or *regexp.Regexp instance. definition should be a function
// which takes special step arguments first and then regular step arguments
// next (with string, int64, *big.Int, and *apd.Decimal
// as valid types) and gocuke.DocString or gocuke.DataTable
// as the last argument if this step uses a doc string or data table respectively.
// Custom step definitions will always take priority of auto-discovered step
// definitions.
func (r *Runner) Step(step interface{}, definition interface{}) *Runner {
	r.topLevelT.Helper()

	exp, ok := step.(*regexp.Regexp)
	if !ok {
		str, ok := step.(string)
		if !ok {
			r.topLevelT.Fatalf("expected step %v to be a string or regex", step)
		}

		var err error
		exp, err = regexp.Compile(str)
		assert.NilError(r.topLevelT, err)
	}

	_ = r.addStepDef(r.topLevelT, exp, reflect.ValueOf(definition))

	return r
}

func (r *Runner) addStepDef(t *testing.T, exp *regexp.Regexp, definition reflect.Value) *stepDef {
	t.Helper()
	defer r.stepDefsMutex.Unlock()

	r.stepDefsMutex.Lock()
	def := r.newStepDefOrHook(t, exp, definition)
	r.stepDefs = append(r.stepDefs, def)
	return def
}

func (r *Runner) newStepDefOrHook(t *testing.T, exp *regexp.Regexp, f reflect.Value) *stepDef {
	t.Helper()

	typ := f.Type()
	if typ.Kind() != reflect.Func {
		t.Fatalf("expected step method, got %s", f)
	}

	funcPtr := f.Pointer()
	rfunc := runtime.FuncForPC(funcPtr)
	file, line := rfunc.FileLine(funcPtr)

	def := &stepDef{
		regex:   exp,
		theFunc: f,
		funcLoc: fmt.Sprintf("%s (%s:%d)", rfunc.Name(), file, line),
	}

	for i := 0; i < typ.NumIn(); i++ {
		stepTyp := typ.In(i)
		getter, ok := r.supportedSpecialArgs[stepTyp]
		if !ok {
			// expect remaining args to be step arguments
			break
		}

		def.specialArgs = append(def.specialArgs, &specialArg{
			typ:      stepTyp,
			getValue: getter,
		})
	}

	if typ.NumOut() != 0 {
		t.Fatalf("expected 0 out parameters for method %+v", f.String())
	}

	return def
}

func (s stepDef) usesRapid() bool {
	for _, arg := range s.specialArgs {
		if arg.typ == rapidTType {
			return true
		}
	}
	return false
}
