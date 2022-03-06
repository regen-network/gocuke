package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	"github.com/cucumber/messages-go/v16"
	"testing"
)

type docRunner struct {
	*Runner
	doc        *messages.GherkinDocument
	astNodeMap map[string]interface{}
}

func newDocRunner(runner *Runner, doc *messages.GherkinDocument) *docRunner {
	astNodeMap := map[string]interface{}{}

	registerSteps := func(steps []*messages.Step) {
		for _, step := range steps {
			astNodeMap[step.Id] = step
		}
	}

	registerBackend := func(background *messages.Background) {
		if background != nil {
			astNodeMap[background.Id] = background
			registerSteps(background.Steps)
		}
	}

	registerScenario := func(scenario *messages.Scenario) {
		if scenario != nil {
			astNodeMap[scenario.Id] = scenario
			registerSteps(scenario.Steps)
			for _, example := range scenario.Examples {
				astNodeMap[example.Id] = example
			}
		}
	}

	for _, child := range doc.Feature.Children {
		registerBackend(child.Background)
		registerScenario(child.Scenario)

		if child.Rule != nil {
			astNodeMap[child.Rule.Id] = child.Rule
			for _, ruleChild := range child.Rule.Children {
				registerBackend(ruleChild.Background)
				registerScenario(ruleChild.Scenario)
			}
		}
	}
	return &docRunner{
		Runner:     runner,
		doc:        doc,
		astNodeMap: astNodeMap,
	}
}

func (r *docRunner) runDoc(t *testing.T) {
	pickles := gherkin.Pickles(*r.doc, r.doc.Uri, r.incr.NewId)
	for _, pickle := range pickles {
		t.Run(pickle.Name, func(t *testing.T) {
			if r.parallel {
				t.Parallel()
			}

			r.runScenario(t, pickle)
		})
	}
}
