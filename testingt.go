package gocuke

import (
	"bytes"
	"fmt"
)

// TestingT is the common subset of testing methods exposed to test suite
// instances and expected by common assertion and mocking libraries.
type TestingT interface {
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

type testingTWrapper struct {
	TestingT
	skipped bool
	errOut  bytes.Buffer
	logOut  bytes.Buffer
}

func (t *testingTWrapper) Error(args ...interface{}) {
	t.Helper()
	t.errOut.WriteString(fmt.Sprint(args...))
	t.TestingT.Error(args...)
}

func (t *testingTWrapper) Errorf(format string, args ...interface{}) {
	t.Helper()
	t.errOut.WriteString(fmt.Sprintf(format, args...))
	t.TestingT.Errorf(format, args...)
}

func (t *testingTWrapper) Fatal(args ...interface{}) {
	t.Helper()
	t.errOut.WriteString(fmt.Sprint(args...))
	t.TestingT.Fatal(args...)
}

func (t *testingTWrapper) Fatalf(format string, args ...interface{}) {
	t.Helper()
	t.errOut.WriteString(fmt.Sprintf(format, args...))
	t.TestingT.Fatalf(format, args...)
}

func (t *testingTWrapper) Log(args ...interface{}) {
	t.Helper()
	t.logOut.WriteString(fmt.Sprint(args...))
	t.TestingT.Log(args...)
}

func (t *testingTWrapper) Logf(format string, args ...interface{}) {
	t.Helper()
	t.logOut.WriteString(fmt.Sprintf(format, args...))
	t.TestingT.Logf(format, args...)
}

func (t *testingTWrapper) Skip(args ...interface{}) {
	t.Helper()
	t.skipped = true
	t.errOut.WriteString(fmt.Sprint(args...))
	t.TestingT.Skip(args...)
}

func (t *testingTWrapper) SkipNow() {
	t.Helper()
	t.skipped = true
	t.TestingT.SkipNow()
}

func (t *testingTWrapper) Skipf(format string, args ...interface{}) {
	t.Helper()
	t.skipped = true
	t.errOut.WriteString(fmt.Sprintf(format, args...))
	t.TestingT.Skipf(format, args...)
}

var _ TestingT = &testingTWrapper{}
