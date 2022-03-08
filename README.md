# gocuke

`gocuke` is a Gherkin-based BDD testing library for golang. It is inspired by
[godog](https://github.com/cucumber/godog) and [gobdd](https://github.com/go-bdd/gobdd).

## Why another golang BDD library?

Sorry guys. I tried godog and gobdd and wanted a specific developer UX that
I couldn't achieve with either. godog was not a good fit for the same reasons
as that gobdd was created (specifically tight integration with `*testing.T`).
Looking at the source code for gobdd, it wasn't that complex but needed to
be updated to a new versions of cucumber/gherkin-go and cucumber/messages-go
which was basically a complete rewrite anyway.

This library offers:
* tight integration with `*testing.T` (use any standard assertion library or mocking framework)
* support for passing context between steps using suites which offers better 
  type safety than other generic context approaches
* auto-discovery of step definitions defined as test suite methods and step
  definition suggestions for minimal configuration
* property-based testing via https://github.com/flyingmutant/rapid
* user-friendly wrapper for data tables
* support for big integers and big decimals (via https://github.com/cockroachdb/apd)
* parallel test execution by default
* full support for all of the latest Gherkin features including rules (via
  the latest cucumber/gherkin-go and cucumber/messages-go)

## Quick Start

### Step 1: Define some Gherkin

In a file `features/simple.feature`:

```gherkin
Feature: simple

  Scenario Outline: eat cukes
    Given I have <x> cukes
    When I eat <y>
    Then I have <z> left

    Examples:
      | x | y | z |
      | 5 | 3 | 2 |
      | 10 | 2 | 8 |
```

### Step 2: Setup the test suite

In a file simple_test.go:

```go
package simple

import (
	"github.com/aaronc/gocuke"
	"testing"
)

func TestMinimal(t *testing.T) {
	gocuke.NewRunner(t, func(t gocuke.TestingT) gocuke.Suite {
		return &suite{TestingT: t}
	}).Run()
}

type suite struct {
	gocuke.TestingT
}
```

When you run the tests, they should fail and suggest that you add these
step definitions:
```go
func (s *suite) IEat(a int64) {
    panic("TODO")
}

func (s *suite) IHaveLeft(a int64) {
    panic("TODO")
}

func (s *suite) IHaveCukes(a int64) {
    panic("TODO")
}
```

Copy these definitions into `simple_test.go`.

### Step 3: Implement Step Definitions

Now implement the step definitions in `simple_test.go`, adding the
variable `cukes int64` to `suite` which tracks state between tests.

**NOTE:** a new `suite` is constructed for every test case so it is safe
to run tests in parallel, which is the default and what is happening
in this example with each of the test cases in the `Scenario Outline`.

```go
type suite struct {
	gocuke.TestingT
	cukes int64
}

func (s *suite) IHaveCukes(a int64) {
	s.cukes = a
}

func (s *suite) IEat(a int64) {
	s.cukes -= a
}

func (s *suite) IHaveLeft(a int64) {
	if a != s.cukes {
		s.Fatalf("expected %d cukes, have %d", a, s.cukes)
	}
}
```

Your tests should now pass!

## Usage Details

### Custom options

Custom `.feature` search paths can be set using the `Runner.WithPath()` method.

Custom step definitions (not auto-discovered on suites) can be added
using the `Runner.Step()` method. The suite must still be the first argument
in all step definitions.

Parallel tests can be disabled using `Runner.NonParallel()`.

### Supported Param Types

`gocuke` supports the following parameter types:
* `string`
* `int64`
* `*big.Int`
* `*apd.Decimal`

`float64` support is not planned because it is lossy!!

### Doc Strings and Data Tables

`gocuke.DocString` or `gocuke.DataTable` should be used as the last argument
in a step definition if the step uses a doc string or data table.

### Property-based testing using Rapid

Property-based tests using https://github.com/flyingmutant/rapid can be
enabled by using `*rapid.T` as the first argument of test methods (after the
suite receiver argument). Property-based test cases will be run as many times
is rapid is configured to run tests.

Example:
```gherkin
Scenario: any int64 value
  Given any int64 string
  When when I convert it to an int64
  Then I get back the original value
```

```go
type suite struct {
  TestingT
  
  x, parsed   int64
  str    string
}

func (s *valuesSuite) AnyInt64String(t *rapid.T) {
	s.x = rapid.Int64().Draw(t, "x").(int64)
	s.str = fmt.Sprintf("%d", s.x)
}

func (s *valuesSuite) WhenIConvertItToAnInt64() {
  s.parsed = toInt64(s, s.str)
}


func (s *suite) IGetBackTheOriginalValue() {
  assert.Equal(s, s.x, s.parsed)
}
```

