package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"pgregory.net/rapid"
	"reflect"
	"testing"
)

// Runner is a test runner.
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

// NewRunner constructs a new Runner with the provided initScenario function.
// initScenario will be called for each test case returning a new suite instance
// for each test case which can be used for sharing state between steps. It
// is expected that the suite will retain a copy of the TestingT instance
// for usage in each step.
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

// Suite is a dummy interface to mark a test suite.
type Suite interface{}

// WithPath specifies glob paths for the runner to look up .feature files.
// The default is `features/*.feature`.
func (r *Runner) WithPath(paths ...string) *Runner {
	r.paths = append(r.paths, paths...)
	return r
}

// NonParallel instructs the runner not to run tests in parallel (which is the default).
func (r *Runner) NonParallel() *Runner {
	r.parallel = false
	return r
}

var rapidTType = reflect.TypeOf(&rapid.T{})
