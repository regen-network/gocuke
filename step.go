package gocuke

import (
	"gotest.tools/v3/assert"
	"reflect"
	"regexp"
)

type stepDef struct {
	exp      *regexp.Regexp
	f        reflect.Value
	hasRapid bool
}

func newStepDef(t TestingT, suiteType reflect.Type, exp *regexp.Regexp, f reflect.Value) *stepDef {
	typ := f.Type()
	if typ.Kind() != reflect.Func {
		t.Fatalf("expected step method, got %s", f)
	}

	if typ.NumIn() < 1 {
		t.Fatalf("expected at least 1 parameter for method %s", f)
	}

	in0 := typ.In(0)
	if in0 != suiteType {
		t.Fatalf("expected at first parameter of method %s to be of type %v", f, in0)
	}

	hasRapid := false
	if typ.NumIn() > 1 {
		in1 := typ.In(1)
		if in1 == rapidTType {
			hasRapid = true
		}
	}

	if typ.NumOut() != 0 {
		t.Fatalf("expected 0 out parameters for method %+v", f.String())
	}

	return &stepDef{
		exp:      exp,
		f:        f,
		hasRapid: hasRapid,
	}
}

// Step can be used to manually register a step with the runner. step should
// be a string or *regexp.Regexp instance. definition should be a function
// which takes the suite as its first argument (usually an instance method),
// parameter arguments next (with string, int64, *big.Int, and *apd.Decimal
// as valid parameter values) and gocuke.DocString or gocuke.DataTable
// as the last argument if this step uses a doc string or data table respectively.
// Custom step definitions will always take priority of auto-discovered step
// definitions.
func (r *Runner) Step(step interface{}, definition interface{}) *Runner {
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

	r.stepDefs = append(r.stepDefs, newStepDef(r.topLevelT, r.suiteType, exp, reflect.ValueOf(definition)))

	return r
}
