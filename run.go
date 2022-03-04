package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	_ "github.com/cucumber/gherkin-go/v19"
	"github.com/cucumber/messages-go/v16"
	_ "github.com/cucumber/messages-go/v16"
	"gotest.tools/v3/assert"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

type Runner struct {
	topLevelT *testing.T
	suiteType reflect.Type
	paths     []string
	parallel  bool
	incr      *messages.Incrementing
	newSuite  func(*testing.T) Suite
	stepDefs  []*stepDef
}

type Suite interface{}

func NewRunner(t *testing.T, newSuite func(*testing.T) Suite) *Runner {
	s := newSuite(t)
	return &Runner{
		topLevelT: t,
		suiteType: reflect.TypeOf(s),
		newSuite:  newSuite,
		incr:      &messages.Incrementing{},
	}
}

func (r *Runner) AddPath(path string) {
	r.paths = append(r.paths, path)
}

func (r *Runner) Step(step string, method interface{}) {
	exp, err := regexp.Compile(step)
	assert.NilError(r.topLevelT, err)

	val := reflect.ValueOf(method)
	typ := val.Type()
	if typ.Kind() != reflect.Func {
		r.topLevelT.Fatalf("expected step method, got %+v", method)
	}

	if typ.NumIn() < 1 {
		r.topLevelT.Fatalf("expected at least 1 parameter for method %+v", method)
	}

	in0 := typ.In(0)
	if in0 != r.suiteType {
		r.topLevelT.Fatalf("expected at first parameter of method %v to be of type %v", method, in0)
	}

	if typ.NumOut() != 0 {
		r.topLevelT.Fatalf("expected 0 out parameters for method %+v", method)
	}

	r.stepDefs = append(r.stepDefs, &stepDef{
		exp: exp,
		f:   val,
	})
}

func (r *Runner) Run() {
	paths := r.paths
	if len(paths) == 0 {
		paths = []string{"features/*.feature"}
	}

	for _, path := range paths {
		files, err := filepath.Glob(path)
		assert.NilError(r.topLevelT, err)

		for _, file := range files {
			r.topLevelT.Run(file, func(t *testing.T) {
				if r.parallel {
					t.Parallel()
				}

				f, err := os.Open(file)
				assert.NilError(t, err)
				defer func() {
					err := f.Close()
					if err != nil {
						panic(err)
					}
				}()

				doc, err := gherkin.ParseGherkinDocument(f, r.incr.NewId)
				assert.NilError(t, err)
				r.runDoc(t, doc)
			})
		}
	}
}

func (r *Runner) runDoc(t *testing.T, doc *messages.GherkinDocument) {
	t.Logf("feature %s", doc.Feature.Name)

	pickles := gherkin.Pickles(*doc, doc.Uri, r.incr.NewId)
	for _, pickle := range pickles {
		t.Run(pickle.Name, func(t *testing.T) {
			t.Logf("pickle %s", pickle.Name)

			if r.parallel {
				t.Parallel()
			}

			s := r.newSuite(t)
			if s, ok := s.(BeforeScenario); ok {
				s.BeforeScenario()
			}

			if s, ok := s.(CleanupScenario); ok {
				defer s.CleanupScenario()
			}

			for _, step := range pickle.Steps {
				if s, ok := s.(SetupStep); ok {
					s.SetupStep()
				}

				r.runStep(t, s, step)

				if s, ok := s.(CleanupStep); ok {
					s.CleanupStep()
				}
			}
		})
	}
}

func (r *Runner) runStep(t *testing.T, s Suite, step *messages.PickleStep) {
	t.Logf("step %s %+v", step.Text, step.Argument)

	// find step
	for _, def := range r.stepDefs {
		matches := def.exp.FindSubmatch([]byte(step.Text))
		if len(matches) == 0 {
			continue
		}

		values := []reflect.Value{reflect.ValueOf(s)}

		for _, _ = range matches[1:] {
			panic("TODO")
		}

		def.f.Call(values)
		return
	}

	t.Fatalf("can't find step definition: %s", step.Text)
}

type SetupStep interface {
	SetupStep()
}

type CleanupStep interface {
	CleanupStep()
}

type BeforeScenario interface {
	BeforeScenario()
}

type CleanupScenario interface {
	CleanupScenario()
}

type stepDef struct {
	exp *regexp.Regexp
	f   reflect.Value
}

type ScenarioContext struct{}

func Run(t *testing.T, runScenario func(*testing.T, *ScenarioContext)) {

}
