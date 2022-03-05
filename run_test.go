package gocuke

import "testing"

func TestRun(t *testing.T) {
	Run(t, Options{}, func(t *testing.T, ctx *ScenarioContext) {
		s := &suite{}
		ctx.Step("I have a data table", s.iHaveADataTable)
		ctx.Step(`some doc string`, s.someDocString)
		ctx.Step(`add (\d+)`, s.add)
		ctx.Step("pass", s.pass)
	})
}

type suite struct{}

func (s *suite) iHaveADataTable(dt DataTable) {}
func (s *suite) someDocString(ds DocString)   {}
func (s *suite) add(x int64, ds DocString)    {}
func (s *suite) pass()                        {}
