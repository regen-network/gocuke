package gocuke

import (
	"testing"

	"gotest.tools/v3/assert"
	"pgregory.net/rapid"
)

func TestHooks(t *testing.T) {
	assert.Assert(t, !longRun)
	assert.Assert(t, !shortRun)
	// test tag expression
	NewRunner(t, &hooksSuite{}).
		Path("features/hooks.feature").
		Tags("@long").
		NonParallel().
		Run()
	assert.Assert(t, longRun)
	assert.Assert(t, !shortRun)

	if open != 0 {
		t.Fatalf("expected 0 open resources, got: %d", open)
	}

	NewRunner(t, &hooksSuite{}).
		Path("features/hooks.feature").
		ShortTags("not @long").
		NonParallel().
		Run()

	assert.Assert(t, longRun)
	assert.Assert(t, shortRun)

	if open != 0 {
		t.Fatalf("expected 0 open resources, got: %d", open)
	}
}

var (
	longRun        = false
	shortRun       = false
	open     int64 = 0
)

type hooksSuite struct {
	TestingT
	numOpenForScenario int64
	numOpenForStep     int64
	numOpenForCleanup  int64
	scenario           Scenario
}

func (s *hooksSuite) IOpenAResource(step Step) {
	assert.Equal(s, "I open a resource", step.Text())
	shortRun = true
	s.numOpenForScenario = 1
	open += s.numOpenForScenario
}

func (s *hooksSuite) IOpenAnyResources(t *rapid.T) {
	longRun = true
	s.numOpenForScenario = rapid.Int64Range(1, 100).AsAny().Draw(t, "numResources").(int64)
	open += s.numOpenForScenario
}

func (s *hooksSuite) IOpenAResourceWithCleanup(step Step) {
	assert.Equal(s, "I open a resource with cleanup", step.Text())
	s.numOpenForScenario = 1
	s.numOpenForCleanup = 1
	open += s.numOpenForScenario + s.numOpenForCleanup
	s.Cleanup(func() {
		open -= s.numOpenForCleanup
	})
}

func (s *hooksSuite) ItIsOpen(step Step) {
	assert.Equal(s, "it is open", step.Text())
	if open < s.numOpenForScenario {
		s.Fatalf("expected at least %d resources to be open", s.numOpenForScenario)
	}
}

func (s *hooksSuite) ExpectScenarioName(a string) {
	assert.Assert(s, s.scenario != nil)
	assert.Equal(s, a, s.scenario.Name())
}

func (s *hooksSuite) ExpectScenarioTag(a string) {
	assert.Equal(s, 1, len(s.scenario.Tags()))
	assert.Equal(s, a, s.scenario.Tags()[0])
}

func (s *hooksSuite) Before(scenario Scenario) {
	s.scenario = scenario
}

func (s *hooksSuite) BeforeStep() {
	if s.numOpenForStep != 0 {
		s.Fatalf("expected step resources to be 0 before step")
	}
	s.numOpenForStep = 1
}

func (s *hooksSuite) AfterStep() {
	if s.numOpenForStep != 1 {
		s.Fatalf("expected step resources to be 1 after step, got: %d", s.numOpenForStep)
	}
	s.numOpenForStep = 0
}

func (s *hooksSuite) After() {
	open -= s.numOpenForScenario
}
