package gocuke_test

import (
	"github.com/aaronc/gocuke"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestSimple(t *testing.T) {
	t.Parallel()
	//gocuke.NewRunner(t, func(t gocuke.TestingT) gocuke.StepDefinitions {
	//	return &simpleSuite{TestingT: t}
	//}).Path("features/simple.feature").Run()
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
	gocuke.NewRunner(t, &customStepsSuite{}).
		Path("features/simple.feature").
		Step(`I have (\d+) cukes`, (*customStepsSuite).A).
		Step(regexp.MustCompile(`I eat (\d+)`), (*customStepsSuite).B).
		Step(`I have (\d+) left`, (*customStepsSuite).C).
		Before((*customStepsSuite).before).
		After((*customStepsSuite).after).
		BeforeStep((*customStepsSuite).beforeStep).
		AfterStep((*customStepsSuite).afterStep).
		NonParallel().
		Run()

	require.Equal(t, 2, beforeCalled)
	require.Equal(t, 2, afterCalled)
	require.Equal(t, 6, beforeStepCalled)
	require.Equal(t, 6, afterStepCalled)
}

var (
	beforeCalled     int
	afterCalled      int
	beforeStepCalled int
	afterStepCalled  int
)

type customStepsSuite struct {
	gocuke.TestingT
	cukes int64
}

func (c *customStepsSuite) before(t gocuke.TestingT) {
	c.TestingT = t
	beforeCalled += 1
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

func (c *customStepsSuite) after() {
	afterCalled += 1
}

func (c *customStepsSuite) beforeStep() {
	beforeStepCalled += 1
}

func (c *customStepsSuite) afterStep() {
	afterStepCalled += 1
}
