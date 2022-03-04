package gocuke

import "testing"

func TestRun(t *testing.T) {
	Run(t, Options{}, func(t *testing.T, ctx *ScenarioContext) {
		s := &suite{}
		ctx.Step("this", s.this)
		ctx.Step(`do (\d+)`, s.doSomething)
		ctx.Step("pass", s.pass)
	})
}

type suite struct{}

func (s *suite) this()             {}
func (s *suite) doSomething(x int) {}
func (s *suite) pass()             {}
