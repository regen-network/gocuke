package gocuke

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
