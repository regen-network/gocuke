package gocuke

import "github.com/aaronc/gocuke/internal/tag"

import "reflect"

// Path specifies glob paths for the runner to look up .feature files.
// The default is `features/*.feature`.
func (r *Runner) Path(paths ...string) *Runner {
	r.paths = append(r.paths, paths...)
	return r
}

// NonParallel instructs the runner not to run tests in parallel (which is the default).
func (r *Runner) NonParallel() *Runner {
	r.parallel = false
	return r
}

func (r *Runner) Before(hook interface{}) *Runner {
	r.addHook(&r.beforeHooks, reflect.ValueOf(hook))
	return r
}

func (r *Runner) After(hook interface{}) *Runner {
	r.addHook(&r.afterHooks, reflect.ValueOf(hook))
	return r
}

func (r *Runner) BeforeStep(hook interface{}) *Runner {
	r.addHook(&r.beforeStepHooks, reflect.ValueOf(hook))
	return r
}

func (r *Runner) AfterStep(hook interface{}) *Runner {
	r.addHook(&r.afterStepHooks, reflect.ValueOf(hook))
	return r
}

func (r *Runner) addHook(hooks *[]*stepDef, f reflect.Value) {
	def := r.newStepDefOrHook(r.topLevelT, nil, f)
	if def.usesRapid() {
		r.suiteUsesRapid = true
	}
	*hooks = append(*hooks, def)
}

// Tags will run only the tests selected by the provided tag expression.
// Tags expressions use the keywords and, or and not and allow expressions
// in parentheses to allow expressions like "(@smoke or @ui) and (not @slow)".
func (r *Runner) Tags(tagExpr string) *Runner {
	var err error
	r.tagExpr, err = tag.ParseExpr(tagExpr)
	if err != nil {
		r.topLevelT.Fatalf("error parsing tag expression %s: %v", tagExpr, err)
	}
	return r
}

// ShortTags specifies which tag expression will be used to select for tests
// when testing.Short() is true. This tag expression will be combined with
// any other tag expression that is applied with Tags() when running short tests.
func (r *Runner) ShortTags(tagExpr string) *Runner {
	var err error
	r.shortTagExpr, err = tag.ParseExpr(tagExpr)
	if err != nil {
		r.topLevelT.Fatalf("error parsing tag expression %s: %v", tagExpr, err)
	}
	return r
}
