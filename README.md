# gocuke :cucumber:

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronc/gocuke.svg)](https://pkg.go.dev/github.com/aaronc/gocuke)

`gocuke` is a Gherkin-based BDD testing library for golang. 

## Features

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

## Why another golang BDD library?

`gocuke` was inspired by
[godog](https://github.com/cucumber/godog) and [gobdd](https://github.com/go-bdd/gobdd).
I tried both of these libraries and wanted a specific developer UX that
I couldn't achieve with either. godog was not a good fit for the same reasons
as that gobdd was created (specifically tight integration with `*testing.T`).
Looking at the source code for gobdd, it needed to
be updated to a new versions of cucumber/gherkin-go and cucumber/messages-go
and significant changes were needed to accommodate this API. So `gocuke` was
written. We are happy to coordinate with the authors
of either of these libraries at some point to align on common goals.

## Quick Start

### Step 1: Write some Gherkin

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
	// a new step definition suite is constructed for every scenario
	gocuke.NewRunner(t, &suite{}).Run()
}

type suite struct {
	// special arguments like TestingT are injected automatically into exported fields
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

### Step Argument Types

`gocuke` supports the following step argument types for arguments captured
from steps:
* `string`
* `int64`
* `*big.Int`
* `*apd.Decimal`

`float64` support is not planned because it is lossy.

### Doc Strings and Data Tables

`gocuke.DocString` or `gocuke.DataTable` should be used as the last argument
in a step definition if the step uses a doc string or data table. `gocuke.DataTable`
provides useful helper methods for working with data tables.

### Special Step Argument Types

The following special argument types are supported:
* `gocuke.TestingT`
* `gocuke.Scenario`
* `gocuke.Step` (will be `nil` when used in a before hook or injected into a suite)
* `*rapid.T` (see Property-based testing using Rapid below)

These can be used in step definitions, hooks, and will be injected into the
suite type if there are exported fields defined with these types.

### Hooks Methods

If the methods `Before`, `After`, `BeforeStep`, or `AfterStep` are defined
on the suite, they will be registered as hooks. `After`  will always be called
even when tests fail. `AfterStep` will always be called whenever a step
started and failed.

It is generally not recommended to over-use hooks. `Before` should primarily be
used for setting up generic resources and `After` should be used for cleaning up
resources. `Given` and `Background` steps should generally be used for setting
up specific test conditions. `BeforeStep` and `AfterStep` should only be used
in very special circumstances.

### Tag Expressions

Cucumber [tag expressions](https://cucumber.io/docs/cucumber/api/#tag-expressions)
can be used for selecting a subset of tests to run. The command-line
option `-gocuke.tags` can be used to specify a subset of tests to run.

The `Runner.Tags()` method can be used to select a set of tags to run in unit
tests. `Runner.ShortTags` method can be used to select a set of tags to

### Custom options

`Runner` has the following methods for setting custom options

* `Path()` sets custom paths. The default is `features/*.feature`.
* `Step()` can be used to add custom steps with special regular expressions.
* `Before()`, `After()`, `BeforeStep()`, or and `AfterStep()` can be used to register custom hooks.
* `Tags` and `ShortTags` can be used with tag expressions as described above.
* `NonParallel()` disables parallel tests.

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
## Roadmap

* [Cucumber `message` based reporting](https://cucumber.io/docs/cucumber/reporting/)
