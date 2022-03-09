package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	"gotest.tools/v3/assert"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Run runs the features registered with the runner.
func (r *Runner) Run() {
	r.topLevelT.Helper()

	paths := r.paths
	if len(paths) == 0 {
		paths = []string{"features/*.feature"}
	}

	for _, path := range paths {
		files, err := filepath.Glob(path)
		assert.NilError(r.topLevelT, err)

		for _, file := range files {
			f, err := os.Open(file)
			assert.NilError(r.topLevelT, err)
			defer func() {
				err := f.Close()
				if err != nil {
					panic(err)
				}
			}()

			doc, err := gherkin.ParseGherkinDocument(f, r.incr.NewId)
			assert.NilError(r.topLevelT, err)
			r.topLevelT.Run(doc.Feature.Name, func(t *testing.T) {
				t.Helper()

				if r.parallel {
					t.Parallel()
				}

				(newDocRunner(r, doc)).runDoc(t)
			})
		}
	}

	if len(r.suggestions) != 0 {
		var suiteTypeName string
		if r.suiteType.Kind() == reflect.Ptr {
			suiteTypeName = "*" + r.suiteType.Elem().Name()
		} else {
			suiteTypeName = r.suiteType.Name()
		}

		suggestionText := "Missing step definitions can be fixed with the following methods:\n"
		for _, sig := range r.suggestions {
			suggestionText += sig.suggestion(suiteTypeName) + "\n\n"
		}
		r.topLevelT.Logf(suggestionText)
	}
}
