package gocuke

import "testing"

func TestRun(t *testing.T) {
	r := NewRunner(t, func(t *testing.T) Suite {
		return &suite{}
	})
	s := &suite{}
	r.Step("pass", s.pass)
	r.Run()
}

type suite struct{}

func (s *suite) pass() {

}
