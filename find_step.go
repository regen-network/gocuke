package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"testing"
)

func (r *Runner) findStep(t *testing.T, step *messages.PickleStep) *stepDef {
	for _, def := range r.stepDefs {
		matches := def.exp.FindSubmatch([]byte(step.Text))
		if len(matches) != 0 {
			return def
		}
	}

	sig := guessMethodSig(step)
	method, ok := r.suiteType.MethodByName(sig.name)
	if ok {
		return newStepDef(t, r.suiteType, sig.regex, method.Func)
	}

	r.suggestions[sig.name] = sig
	t.Errorf("can't find step definition for: %s", step.Text)

	return nil
}

//func (r *Runner) findStep(t *testing.T, step *messages.PickleStep) stepClosure {
//	for _, def := range r.stepDefs {
//		matches := def.exp.FindSubmatch([]byte(step.Text))
//		if len(matches) == 0 {
//			continue
//		}
//
//		matches = matches[1:]
//		expectedIn := len(matches)
//		typ := def.f.Type()
//
//		hasPickleArg := step.Argument != nil
//		if hasPickleArg {
//			expectedIn += 1
//		}
//
//		if expectedIn != typ.NumIn() {
//			t.Fatalf("expected %d in parameters for function %+v", expectedIn, def.f)
//		}
//
//		values := make([]reflect.Value, expectedIn)
//		for i, match := range matches {
//			values[i] = convertParamValue(t, string(match), typ.In(i))
//		}
//
//		// pickleArg goes last
//		if hasPickleArg {
//			i := expectedIn - 1
//			typ := typ.In(i)
//			// only one of DataTable or DocString is valid
//			if typ == dataTableType {
//				if step.Argument.DataTable == nil {
//					t.Fatalf("expected non-nil DataTable")
//				}
//
//				dataTable := DataTable{
//					t:     t,
//					table: step.Argument.DataTable,
//				}
//				values[i] = reflect.ValueOf(dataTable)
//			} else if typ == docStringType {
//				if step.Argument.DocString == nil {
//					t.Fatalf("expected non-nil DocString")
//				}
//
//				docString := DocString{
//					MediaType: step.Argument.DocString.MediaType,
//					Content:   step.Argument.DocString.Content,
//				}
//				values[i] = reflect.ValueOf(docString)
//			} else {
//				t.Fatalf("unexpected parameter type %v", typ)
//			}
//		}
//
//		def.f.Call(values)
//		return
//	}
//}
//
//type stepClosure struct {
//	*stepDef
//	values []reflect.Value
//}
