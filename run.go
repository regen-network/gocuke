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
	"strconv"
	"testing"
)

type runner struct {
	opts          Options
	setupScenario func(*testing.T, *ScenarioContext)
	topLevelT     *testing.T
	incr          *messages.Incrementing
}

type Options struct {
	Paths    []string
	Parallel bool
}

func (r *runner) run() {
	paths := r.opts.Paths
	if len(paths) == 0 {
		paths = []string{"features/*.feature"}
	}

	for _, path := range paths {
		files, err := filepath.Glob(path)
		assert.NilError(r.topLevelT, err)

		for _, file := range files {
			r.topLevelT.Run(file, func(t *testing.T) {
				if r.opts.Parallel {
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

func (r *runner) runDoc(t *testing.T, doc *messages.GherkinDocument) {
	t.Logf("feature %s", doc.Feature.Name)

	pickles := gherkin.Pickles(*doc, doc.Uri, r.incr.NewId)
	for _, pickle := range pickles {
		t.Run(pickle.Name, func(t *testing.T) {
			t.Logf("pickle %s", pickle.Name)

			if r.opts.Parallel {
				t.Parallel()
			}

			scenarioCtx := &ScenarioContext{
				stepDefs: nil,
				t:        t,
			}

			r.setupScenario(t, scenarioCtx)

			for _, step := range pickle.Steps {
				r.runStep(t, scenarioCtx, step)
			}
		})
	}
}

func (r *runner) runStep(t *testing.T, ctx *ScenarioContext, step *messages.PickleStep) {
	t.Logf("step %s %+v", step.Text, step.Argument)

	// find step
	for _, def := range ctx.stepDefs {
		matches := def.exp.FindSubmatch([]byte(step.Text))
		if len(matches) == 0 {
			continue
		}

		matches = matches[1:]
		n := len(matches)
		typ := def.f.Type()
		if n != typ.NumIn() {
			t.Fatalf("expected %d in parameters for function %+v", n, def.f)
		}

		values := make([]reflect.Value, n)
		for i, match := range matches {
			kind := typ.In(i).Kind()
			switch kind {
			case reflect.Int:
				x, err := strconv.Atoi(string(match))
				assert.NilError(t, err)
				values[i] = reflect.ValueOf(x)
			case reflect.Uint:
				x, err := strconv.Atoi(string(match))
				assert.NilError(t, err)
				values[i] = reflect.ValueOf(uint(x))
			case reflect.String:
				values[i] = reflect.ValueOf(string(match))
			default:
				t.Fatalf("unexpected parameter kind %s", kind)
			}
		}

		def.f.Call(values)
		return
	}

	t.Fatalf("can't find step definition: %s", step.Text)
}

type stepDef struct {
	exp *regexp.Regexp
	f   reflect.Value
}

func (r *ScenarioContext) Step(step string, fn interface{}) {
	exp, err := regexp.Compile(step)
	assert.NilError(r.t, err)

	val := reflect.ValueOf(fn)
	typ := val.Type()
	if typ.Kind() != reflect.Func {
		r.t.Fatalf("expected step fn, got %+v", fn)
	}

	if typ.NumOut() != 0 {
		r.t.Fatalf("expected 0 out parameters for fn %+v", fn)
	}

	r.stepDefs = append(r.stepDefs, &stepDef{
		exp: exp,
		f:   val,
	})
}

type ScenarioContext struct {
	stepDefs []*stepDef
	t        *testing.T
}

func Run(t *testing.T, opts Options, setupScenario func(t *testing.T, ctx *ScenarioContext)) {
	r := &runner{
		opts:          opts,
		setupScenario: setupScenario,
		topLevelT:     t,
		incr:          &messages.Incrementing{},
	}
	r.run()
}
