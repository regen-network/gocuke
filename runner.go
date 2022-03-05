package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"pgregory.net/rapid"
	"reflect"
	"testing"
)

type Runner struct {
	topLevelT    *testing.T
	suiteType    reflect.Type
	initScenario func(t TestingT) Suite
	incr         *messages.Incrementing
	paths        []string
	parallel     bool
	stepDefs     []*stepDef
	suggestions  map[string]methodSig
}

func NewRunner(t *testing.T, initScenario func(t TestingT) Suite) *Runner {
	s := initScenario(t)
	return &Runner{
		topLevelT:    t,
		suiteType:    reflect.TypeOf(s),
		initScenario: initScenario,
		incr:         &messages.Incrementing{},
		suggestions:  map[string]methodSig{},
	}
}

type Suite interface{}

func (r *Runner) WithPath(paths ...string) *Runner {
	r.paths = append(r.paths, paths...)
	return r
}

func (r *Runner) NonParallel() *Runner {
	r.parallel = false
	return r
}

var rapidTType = reflect.TypeOf(&rapid.T{})
