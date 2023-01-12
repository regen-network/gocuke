package gocuke

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	gherkin "github.com/cucumber/gherkin/go/v26"
	"gotest.tools/v3/assert"
)

// Run runs the features registered with the runner.
func (r *Runner) Run() {
	r.topLevelT.Helper()

	paths := r.paths
	if len(paths) == 0 {
		paths = []string{"features/*.feature"}
	}

	haveTests := false

	for _, path := range paths {
		var files []string
		// use glob paths if we have a * in the path
		// if we don't have a glob just check the path directly
		// not doing this allows mis-spellings in exact paths to be skipped silently
		if strings.Contains(path, "*") {
			var err error
			files, err = filepath.Glob(path)
			assert.NilError(r.topLevelT, err)
		} else {
			files = []string{path}
		}

		for _, file := range files {
			f, err := os.Open(file)
			assert.NilError(r.topLevelT, err)
			defer func() {
				err := f.Close()
				if err != nil {
					panic(err)
				}
			}()

			haveTests = true

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
			suggestionText += sig.methodSuggestion(suiteTypeName) + "\n\n"
		}

		suggestionText += "Steps can be manually registered with the runner for customization using this code:\n"
		for _, sig := range r.suggestions {
			suggestionText += "  " + sig.stepSuggestion(suiteTypeName) + ".\n"
		}
		suggestionText += "\n\n"
		suggestionText += "See https://github.com/regen-network/gocuke for further customization options."

		r.topLevelT.Logf(suggestionText)
	}

	if !haveTests {
		r.topLevelT.Fatalf("no tests found in paths: %v", r.paths)
	}
}
