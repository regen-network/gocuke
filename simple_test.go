package gocuke_test

import (
	"testing"

	"github.com/regen-network/gocuke"
)

func TestSimple(t *testing.T) {
	gocuke.NewRunner(t, &simpleSuite{}).Path("features/simple.feature").Run()
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

// test if a struct that doesn't use a pointer and a global var
func TestSimpleNonPointer(t *testing.T) {
	gocuke.NewRunner(t, simpleSuiteNP{}).Path("features/simple.feature").Run()
}

var globalCukes int64

type simpleSuiteNP struct {
	gocuke.TestingT
}

func (s simpleSuiteNP) IHaveCukes(a int64) {
	globalCukes = a
}

func (s simpleSuiteNP) IEat(a int64) {
	globalCukes -= a
}

func (s simpleSuiteNP) IHaveLeft(a int64) {
	if a != globalCukes {
		s.Fatalf("expected %d cukes, have %d", a, globalCukes)
	}
}
