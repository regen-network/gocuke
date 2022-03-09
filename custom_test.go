package gocuke_test

import (
	"github.com/regen-network/gocuke"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestCustomSteps(t *testing.T) {
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

func (c customStepsSuite) before() {
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

func (c customStepsSuite) after() {
	afterCalled += 1
}

func (c customStepsSuite) beforeStep() {
	beforeStepCalled += 1
}

func (c customStepsSuite) afterStep() {
	afterStepCalled += 1
}
