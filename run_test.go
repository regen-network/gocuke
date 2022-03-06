package gocuke_test

import (
	"gocuke"
	"testing"
)

func TestRun(t *testing.T) {
	gocuke.NewRunner(t, func(t gocuke.TestingT) gocuke.Suite {
		return &suite{t}
	}).WithPath("features/test1.feature").Run()
}

type suite struct {
	gocuke.TestingT
}

func (s *suite) IHaveADataTable(dt gocuke.DataTable) {}
func (s *suite) SomeDocString(ds gocuke.DocString)   {}
func (s *suite) Add(x int64, ds gocuke.DocString)    {}
func (s *suite) Pass()                               {}

func (s *suite) IShouldHaveCucumbers(a int64) {
}

func (s *suite) ThereAreCucumbers(a int64) {
}

func (s *suite) IEatCucumbers(a int64) {
}
