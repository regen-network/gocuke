package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	"github.com/cucumber/messages-go/v16"
	"testing"
)

func (r *runner) runDoc(t *testing.T, doc *messages.GherkinDocument) {
	t.Logf("feature %s", doc.Feature.Name)

	pickles := gherkin.Pickles(*doc, doc.Uri, r.incr.NewId)
	for _, pickle := range pickles {
		t.Run(pickle.Name, func(t *testing.T) {
			t.Logf("pickle %s", pickle.Name)

			if r.opts.Parallel {
				t.Parallel()
			}

			scenarioCtx := &ScenarioContext{
				stepDefs: nil,
				t:        t,
				pickle:   pickle,
			}

			r.setupScenario(t, scenarioCtx)

			for _, step := range pickle.Steps {
				r.runStep(t, scenarioCtx, step)
			}
		})
	}
}
