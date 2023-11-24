package gocuke

import (
	"reflect"

	messages "github.com/cucumber/messages/go/v21"
)

func (r *scenarioRunner) runHook(def *stepDef) {
	r.t.Helper()

	typ := def.theFunc.Type()
	expectedIn := len(def.specialArgs)
	if expectedIn != typ.NumIn() {
		r.t.Fatalf("expected %d in parameter(s) for function %+v, got %d", expectedIn, def.funcLoc, typ.NumIn())
	}

	values := make([]reflect.Value, expectedIn)

	for i, arg := range def.specialArgs {
		values[i] = reflect.ValueOf(arg.getValue(r))
	}

	def.theFunc.Call(values)
}

func (r *scenarioRunner) runStep(step *messages.PickleStep, def *stepDef) {
	r.t.Helper()

	r.step = step

	for _, hook := range r.beforeStepHooks {
		r.runHook(hook)
	}

	for _, hook := range r.afterStepHooks {
		defer r.runHook(hook)
	}

	matches := def.regex.FindSubmatch([]byte(step.Text))
	if len(matches) == 0 {
		r.t.Fatalf("internal error: no matches found when matching %s against %s", def.regex.String(), step.Text)
	}

	matches = matches[1:]
	numSpecialArgs := len(def.specialArgs)
	expectedIn := len(matches) + numSpecialArgs
	typ := def.theFunc.Type()

	hasPickleArg := step.Argument != nil
	if hasPickleArg {
		expectedIn += 1
	}

	if expectedIn != typ.NumIn() {
		r.t.Fatalf("expected %d in parameter(s) for function %+v, got %d", expectedIn, def.funcLoc, typ.NumIn())
	}

	values := make([]reflect.Value, expectedIn)

	for i, arg := range def.specialArgs {
		values[i] = reflect.ValueOf(arg.getValue(r))
	}

	for i, match := range matches {
		values[i+numSpecialArgs] = convertParamValue(r.t, string(match), typ.In(i+numSpecialArgs), def.funcLoc)
	}

	// pickleArg goes last
	if hasPickleArg {
		i := expectedIn - 1
		pickleArgType := typ.In(i)
		// only one of DataTable or DocString is valid
		if pickleArgType == dataTableType {
			if step.Argument.DataTable == nil {
				r.t.Fatalf("expected non-nil DataTable")
			}

			dataTable := DataTable{
				t:     r.t,
				table: step.Argument.DataTable,
			}
			values[i] = reflect.ValueOf(dataTable)
		} else if pickleArgType == docStringType {
			if step.Argument.DocString == nil {
				r.t.Fatalf("expected non-nil DocString")
			}

			docString := DocString{
				MediaType: step.Argument.DocString.MediaType,
				Content:   step.Argument.DocString.Content,
			}
			values[i] = reflect.ValueOf(docString)
		} else {
			r.t.Fatalf("unexpected parameter type %v in function %s", pickleArgType, def.funcLoc)
		}
	}

	if r.verbose {
		r.t.Logf("Step: %s", step.Text)
	}

	def.theFunc.Call(values)
}
