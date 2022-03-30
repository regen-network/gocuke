package gocuke

import (
	"github.com/cucumber/common/gherkin/go/v22"
	"github.com/cucumber/common/messages/go/v17"
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
