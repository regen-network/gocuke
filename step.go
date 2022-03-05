package gocuke

import (
	"gotest.tools/v3/assert"
	"reflect"
	"regexp"
)

type stepDef struct {
	exp *regexp.Regexp
	f   reflect.Value
}

func (r *ScenarioContext) Step(step string, fn interface{}) {
	exp, err := regexp.Compile(step)
	assert.NilError(r.t, err)

	val := reflect.ValueOf(fn)
	r.addStep(exp, val)
}

func (r *ScenarioContext) addStep(exp *regexp.Regexp, val reflect.Value) {
	typ := val.Type()
	if typ.Kind() != reflect.Func {
		r.t.Fatalf("expected step fn, got %+v", val)
	}

	if typ.NumOut() != 0 {
		r.t.Fatalf("expected 0 out parameters for fn %+v", val)
	}

	r.stepDefs = append(r.stepDefs, &stepDef{
		exp: exp,
		f:   val,
	})
}
