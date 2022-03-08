package gocuke

import (
	"pgregory.net/rapid"
	"testing"
)

func TestCleanup(t *testing.T) {
	NewRunner(t, func(t TestingT) Suite {
		return &cleanupSuite{TestingT: t}
	}).WithPath("features/cleanup.feature").Run()

	if open != 0 {
		t.Fatalf("expected resource to be closed")
	}
}

var open int64 = 0

type cleanupSuite struct {
	TestingT
	numOpen int64
}

func (s *cleanupSuite) IOpenAResource() {
	s.numOpen = 1
	open += s.numOpen
}

func (s *cleanupSuite) IOpenAnyResources(t *rapid.T) {
	s.numOpen = rapid.Int64Range(1, 100).Draw(t, "numResources").(int64)
	open += s.numOpen
}

func (s *cleanupSuite) ItIsOpen() {
	if open < s.numOpen {
		s.Fatalf("expected at least %d resources to be open", s.numOpen)
	}
}

func (s *cleanupSuite) Cleanup() {
	open -= s.numOpen
}
