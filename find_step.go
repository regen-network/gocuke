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