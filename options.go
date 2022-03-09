package gocuke

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

// Before registers a before hook function which can take special step arguments.
func (r *Runner) Before(hook interface{}) *Runner {
	r.addHook(&r.beforeHooks, reflect.ValueOf(hook))
	return r
}

// After registers an after hook function which can take special step arguments.
// This hook will be called even when the test fails.
func (r *Runner) After(hook interface{}) *Runner {
	r.addHook(&r.afterHooks, reflect.ValueOf(hook))
	return r
}

// BeforeStep registers a before step hook function which can take special step arguments.
func (r *Runner) BeforeStep(hook interface{}) *Runner {
	r.addHook(&r.beforeStepHooks, reflect.ValueOf(hook))
	return r
}

// AfterStep registers an after step hook function which can take special step arguments.
// This hook will be called even when the test fails.
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
