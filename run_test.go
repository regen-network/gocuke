package gocuke_test

import (
	"gocuke"
	"testing"
)

func TestRun(t *testing.T) {
	gocuke.NewRunner(t, func(t gocuke.TestingT) gocuke.Suite {
		return &suite{t}
	}).Run()
}

type suite struct {
	gocuke.TestingT
}

func (s *suite) IHaveADataTable(dt gocuke.DataTable) {}
func (s *suite) SomeDocString(ds gocuke.DocString)   {}
func (s *suite) Add(x int64, ds gocuke.DocString)    {}
func (s *suite) Pass()                               {}
