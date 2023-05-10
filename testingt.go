package gocuke

import "pgregory.net/rapid"

// TestingT is the common subset of testing methods exposed to test suite
// instances and expected by common assertion and mocking libraries.
type TestingT interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Helper()
}

// rapidT is a wrapper around `*rapid.T` that stubs missing `TestingT`
// interface members (e.g. `Cleanup()`).
type rapidT struct {
	*rapid.T
}

func (rt *rapidT) Cleanup(fn func()) {
	rt.Log("WARNING: cleanup called on `*rapid.T`")
}
