package gocuke

import "testing"

func TestRun(t *testing.T) {
	Run(t, Options{}, func(t *testing.T, ctx *ScenarioContext) {
		ctx.StepSuite(&suite{t: t})
	})
}

type suite struct {
	t *testing.T
}

func (s *suite) IHaveADataTable(dt DataTable) {}
func (s *suite) SomeDocString(ds DocString)   {}
func (s *suite) Add(x int64, ds DocString)    {}
func (s *suite) Pass()                        {}
