package gocuke

import (
	"testing"

	messages "github.com/cucumber/messages/go/v22"
)

func (r *Runner) findStep(t *testing.T, step *messages.PickleStep) *stepDef {
	t.Helper()

	r.stepDefsMutex.RLock()
	def := findSubmatch(t, step.Text, r.stepDefs)
	r.stepDefsMutex.RUnlock()
	if def != nil {
		return def
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

func findSubmatch(t *testing.T, stepText string, stepDefs []*stepDef) *stepDef {
	t.Helper()

	for _, def := range stepDefs {
		matches := def.regex.FindSubmatch([]byte(stepText))
		if len(matches) != 0 {
			return def
		}
	}

	return nil
}