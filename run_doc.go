package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	"github.com/cucumber/messages-go/v16"
	"testing"
)

func (r *Runner) runDoc(t *testing.T, doc *messages.GherkinDocument) {
	pickles := gherkin.Pickles(*doc, doc.Uri, r.incr.NewId)
	for _, pickle := range pickles {
		t.Run(pickle.Name, func(t *testing.T) {
			if r.parallel {
				t.Parallel()
			}

			r.runScenario(t, pickle)
		})
	}
}
