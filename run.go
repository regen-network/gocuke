package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	"gotest.tools/v3/assert"
	"os"
	"path/filepath"
	"testing"
)

func (r *Runner) Run() {
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
			r.topLevelT.Run(file, func(t *testing.T) {
				if r.parallel {
					t.Parallel()
				}

				r.runDoc(t, doc)
			})
		}
	}

	//defer func() {
	if len(r.suggestions) != 0 {
		suggestionText := "Missing definitions can be fixed with the following methods:\n"
		for _, sig := range r.suggestions {
			suggestionText += sig.suggestion(r.suiteType) + "\n\n"
		}
		r.topLevelT.Logf(suggestionText)
	}
	//}()
}
