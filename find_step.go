package gocuke

import (
	"testing"

	messages "github.com/cucumber/messages/go/v21"
)

func (r *Runner) findStep(t *testing.T, step *messages.PickleStep) *stepDef {
	t.Helper()

	for _, def := range r.stepDefs {
		matches := def.regex.FindSubmatch([]byte(step.Text))
		if len(matches) != 0 {
			return def
		}
	}

	sig := guessMethodSig(step)
	method, ok := r.suiteType.MethodByName(sig.name)
	if ok {
		return r.addStepDef(t, sig.regex, method.Func)
	}

	if !r.haveSuggestion[sig.name] {
		r.haveSuggestion[sig.name] = true
		r.suggestions = append(r.suggestions, sig)
	}
	t.Errorf("can't find step definition for: %s", step.Text)

	return nil
}
