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
	typ := val.Type()
	if typ.Kind() != reflect.Func {
		r.t.Fatalf("expected step fn, got %+v", fn)
	}

	if typ.NumOut() != 0 {
		r.t.Fatalf("expected 0 out parameters for fn %+v", fn)
	}

	r.stepDefs = append(r.stepDefs, &stepDef{
		exp: exp,
		f:   val,
	})
}
