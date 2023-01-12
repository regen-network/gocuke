package gocuke

import (
	"testing"

	gherkin "github.com/cucumber/gherkin/go/v26"
	messages "github.com/cucumber/messages/go/v21"
)

type docRunner struct {
	*Runner
	doc *messages.GherkinDocument
}

func newDocRunner(runner *Runner, doc *messages.GherkinDocument) *docRunner {
	return &docRunner{
		Runner: runner,
		doc:    doc,
	}
}

func (r *docRunner) runDoc(t *testing.T) {
	t.Helper()

	pickles := gherkin.Pickles(*r.doc, r.doc.Uri, r.incr.NewId)
	for _, pickle := range pickles {
		t.Run(pickle.Name, func(t *testing.T) {
			t.Helper()

			if r.parallel {
				t.Parallel()
			}

			r.runScenario(t, pickle)
		})
	}
}
