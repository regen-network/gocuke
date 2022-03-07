package gocuke

// TestingT is the common subset of testing methods exposed to test suite
// instances and expected by common assertion and mocking libraries.
type TestingT interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Log(args ...any)
	Logf(format string, args ...any)
	Skip(args ...any)
	SkipNow()
	Skipf(format string, args ...any)
	Helper()
}
