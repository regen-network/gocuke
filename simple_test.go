package gocuke_test

import (
	"gocuke"
	"testing"
)

func TestMinimal(t *testing.T) {
	gocuke.NewRunner(t, func(t gocuke.TestingT) gocuke.Suite {
		return &suite{TestingT: t}
	}).WithPath("features/simple.feature").Run()
}

type suite struct {
	gocuke.TestingT
	cukes int64
}

func (s *suite) IHaveCukes(a int64) {
	s.cukes = a
}

func (s *suite) IEat(a int64) {
	s.cukes -= a
}

func (s *suite) IHaveLeft(a int64) {
	if a != s.cukes {
		s.Fatalf("expected %d cukes, have %d", a, s.cukes)
	}
}
