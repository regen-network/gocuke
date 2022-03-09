package gocuke

import (
	"gotest.tools/v3/assert"
	"pgregory.net/rapid"
	"testing"
)

func TestHooks(t *testing.T) {
	NewRunner(t, &hooksSuite{}).Path("features/hooks.feature").Run()

	if open != 0 {
		t.Fatalf("expected resource to be closed")
	}
}

var open int64 = 0

type hooksSuite struct {
	TestingT
	numOpenForScenario int64
	numOpenForStep     int64
	scenario           Scenario
}

func (s *hooksSuite) IOpenAResource() {
	s.numOpenForScenario = 1
	open += s.numOpenForScenario
}

func (s *hooksSuite) IOpenAnyResources(t *rapid.T) {
	s.numOpenForScenario = rapid.Int64Range(1, 100).Draw(t, "numResources").(int64)
	open += s.numOpenForScenario
}

func (s *hooksSuite) ItIsOpen() {
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
		s.Fatalf("expected step resources to be 1 before step")
	}
	s.numOpenForStep = 0
}

func (s *hooksSuite) After() {
	open -= s.numOpenForScenario
}
