package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"pgregory.net/rapid"
	"reflect"
)

func (r *Runner) runStep(t TestingT, step *messages.PickleStep, def *stepDef, s Suite) {
	matches := def.exp.FindSubmatch([]byte(step.Text))
	matches = matches[1:]
	expectedIn := len(matches) + 1
	typ := def.f.Type()

	if def.hasRapid {
		expectedIn += 1
	}

	hasPickleArg := step.Argument != nil
	if hasPickleArg {
		expectedIn += 1
	}

	if expectedIn != typ.NumIn() {
		t.Fatalf("expected %d in parameter(s) for function %+v, got %d", expectedIn, def.f.String(), typ.NumIn())
	}

	values := make([]reflect.Value, expectedIn)
	values[0] = reflect.ValueOf(s)

	offset := 1
	if def.hasRapid {
		values[offset] = reflect.ValueOf(t.(*rapid.T))
		offset++
	}

	for i, match := range matches {
		values[i+offset] = convertParamValue(t, string(match), typ.In(i+offset))
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
}
