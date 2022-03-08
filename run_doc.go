package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	"github.com/cucumber/messages-go/v16"
	"testing"
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
