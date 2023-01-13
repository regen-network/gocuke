package gocuke

import (
	"testing"

	messages "github.com/cucumber/messages/go/v21"
)

func (r *Runner) findStep(t *testing.T, step *messages.PickleStep) *stepDef {
	t.Helper()

	for _, def := range r.stepDefs {
		args, err := def.expr.Match(step.Text)
		if err == nil && len(args) != 0 {
			return def
		}
	}

	sig := guessMethodSig(step)
	method, ok := r.suiteType.MethodByName(sig.name)
	if ok {
		return r.addStepDef(t, sig.regex, method.Func)
	}

	r.suggestions[sig.name] = sig
	t.Errorf("can't find step definition for: %s", step.Text)

	return nil
}
