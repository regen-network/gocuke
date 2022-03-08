package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"gotest.tools/v3/assert"
	"testing"
)

func TestGuessMethodSig(t *testing.T) {
	NewRunner(t, func(t TestingT) StepDefinitions {
		return &guessSuite{TestingT: t}
	}).Path("features/guess.feature").Run()
}

type guessSuite struct {
	TestingT

	step         string
	sig          methodSig
	matches      []string
	hasDocString bool
	hasDataTable bool
}

func (s *guessSuite) TheStep(a DocString) {
	s.step = a.Content
}

func (s *guessSuite) WeGuessTheStepDefinition() {
	ps := &messages.PickleStep{Text: s.step}
	if s.hasDocString {
		ps.Argument = &messages.PickleStepArgument{DocString: &messages.PickleDocString{}}
	}
	if s.hasDataTable {
		ps.Argument = &messages.PickleStepArgument{DataTable: &messages.PickleTable{}}
	}
	s.sig = guessMethodSig(ps)
}

func (s *guessSuite) WeGetTheMethodSignature(a DocString) {
	assert.Equal(s, s.sig.methodSig(), a.Content)
}

func (s *guessSuite) WeMatchTheStep() {
	s.matches = s.sig.regex.FindStringSubmatch(s.step)[1:]
}

func (s *guessSuite) WeGetTheValues(a DataTable) {
	assert.Equal(s, a.NumCols(), len(s.matches))
	for i, match := range s.matches {
		assert.Equal(s, a.Cell(0, i).String(), match)
	}
}

func (s *guessSuite) WithADocString() {
	s.hasDocString = true
}

func (s *guessSuite) WithADataTable() {
	s.hasDataTable = true
}
