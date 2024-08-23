package gocuke_test

import (
	"testing"

	"github.com/regen-network/gocuke"
)

func TestSimple(t *testing.T) {
	gocuke.NewRunner(t, &simpleSuite{}).Path("examples/simple/simple.feature").Run()
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
	gocuke.
		NewRunner(t, simpleSuiteNP{}).
		Path("examples/simple/simple.feature").
		NonParallel().
		Run()
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

// test a struct using a different interface compatible with gocuke.TestingT
func TestSimpleCompat(t *testing.T) {
	gocuke.NewRunner(t, &simpleSuiteCompat{}).Path("examples/simple/simple.feature").Run()
}

type TestingTCompat interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Helper()

	// Not included in gocuke.TestingT
	Name() string
}

type simpleSuiteCompat struct {
	TestingTCompat
	cukes int64
}

func (s *simpleSuiteCompat) IHaveCukes(a int64) {
	// These calls to s.LogF fail if s.TestingTCompat is nil
	s.Logf("I have %d cukes", a)
	s.cukes = a
}

func (s *simpleSuiteCompat) IEat(a int64) {
	s.Logf("I eat %d", a)
	s.cukes -= a
}

func (s *simpleSuiteCompat) IHaveLeft(a int64) {
	s.Logf("I have %d left?", a)
	if a != s.cukes {
		s.Fatalf("expected %d cukes, have %d", a, s.cukes)
	}
}
