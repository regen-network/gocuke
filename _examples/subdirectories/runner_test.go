package main

import (
	"testing"

	"github.com/regen-network/gocuke"
)

// should recursively load scenarios from features directory
func TestSimple(t *testing.T) {
	gocuke.NewRunner(t, &simpleSuite{}).Run()
}

type simpleSuite struct {
	gocuke.TestingT
	cukes int64
}

func (s *simpleSuite) IHaveCukes(a int64) {
	s.cukes = a
}

func (s *simpleSuite) IEat(a int64) {
	s.cukes -= a
}

func (s *simpleSuite) IHaveLeft(a int64) {
	if a != s.cukes {
		s.Fatalf("expected %d cukes, have %d", a, s.cukes)
	}
}
