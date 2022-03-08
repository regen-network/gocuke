package gocuke_test

import (
	"github.com/aaronc/gocuke"
	"regexp"
	"testing"
)

func TestSimple(t *testing.T) {
	t.Parallel()
	gocuke.NewRunner(t, func(t gocuke.TestingT) gocuke.Suite {
		return &simpleSuite{TestingT: t}
	}).Path("features/simple.feature").Run()
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

func TestCustomSteps(t *testing.T) {
	t.Parallel()
	gocuke.NewRunner(t, func(t gocuke.TestingT) gocuke.Suite {
		return &customStepsSuite{TestingT: t}
	}).
		Path("features/simple.feature").
		Step(`I have (\d+) cukes`, (*customStepsSuite).A).
		Step(regexp.MustCompile(`I eat (\d+)`), (*customStepsSuite).B).
		Step(`I have (\d+) left`, (*customStepsSuite).C).
		Run()
}

type customStepsSuite struct {
	gocuke.TestingT
	cukes int64
}

func (s *customStepsSuite) A(a int64) {
	s.cukes = a
}

func (s *customStepsSuite) B(a int64) {
	s.cukes -= a
}

func (s *customStepsSuite) C(a int64) {
	if a != s.cukes {
		s.Fatalf("expected %d cukes, have %d", a, s.cukes)
	}
}
