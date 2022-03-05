package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"gotest.tools/v3/assert"
	"reflect"
	"regexp"
	"testing"
)

func (r *runner) runStep(t *testing.T, ctx *ScenarioContext, step *messages.PickleStep) {
	t.Logf("step %s", step.Text)
	for _, def := range ctx.stepDefs {
		matches := def.exp.FindSubmatch([]byte(step.Text))
		if len(matches) == 0 {
			continue
		}

		matches = matches[1:]
		expectedIn := len(matches)
		typ := def.f.Type()

		hasPickleArg := step.Argument != nil
		if hasPickleArg {
			expectedIn += 1
		}

		if expectedIn != typ.NumIn() {
			t.Fatalf("expected %d in parameters for function %+v", expectedIn, def.f)
		}

		values := make([]reflect.Value, expectedIn)
		for i, match := range matches {
			values[i] = convertParamValue(t, string(match), typ.In(i))
		}

		// pickleArg goes last
		if hasPickleArg {
			i := expectedIn - 1
			typ := typ.In(i)
			// only one of DataTable or DocString is valid
			if typ == dataTableType {
				if step.Argument.DataTable == nil {
					t.Fatalf("expected non-nil DataTable")
				}

				dataTable := DataTable{
					t:     t,
					table: step.Argument.DataTable,
				}
				values[i] = reflect.ValueOf(dataTable)
			} else if typ == docStringType {
				if step.Argument.DocString == nil {
					t.Fatalf("expected non-nil DocString")
				}

				docString := DocString{
					MediaType: step.Argument.DocString.MediaType,
					Content:   step.Argument.DocString.Content,
				}
				values[i] = reflect.ValueOf(docString)
			} else {
				t.Fatalf("unexpected parameter type %v", typ)
			}
		}

		def.f.Call(values)
		return
	}

	t.Fatalf("can't find step definition: %s", step.Text)
}

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
