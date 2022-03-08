package gocuke

import (
	"github.com/aaronc/gocuke/internal/tag"
	"gotest.tools/v3/assert"
	"strings"
	"testing"
)

func TestTags(t *testing.T) {
	NewRunner(t, func(t TestingT) Suite {
		return &tagsSuite{TestingT: t}
	}).Path("features/tags.feature").Run()
}

type tagsSuite struct {
	TestingT
	expr *tag.Expr
	tags []string
}

var eatSomeCukes = false

func (s *tagsSuite) IEatSomeCukes() {
	eatSomeCukes = true
}

var eatOtherCukes = false

func (s *tagsSuite) IEatOtherCukes() {
	eatOtherCukes = true
}

func (s *tagsSuite) TheTagExpression(a DocString) {
	var err error
	s.expr, err = tag.ParseExpr(a.Content)
	assert.NilError(s, err, a.Content)
}

func (s *tagsSuite) IMatch(a string) {
	s.tags = strings.Split(a, " ")
}

func (s *tagsSuite) TheResultIs(a string) {
	res := s.expr.Match(tag.NewTags(s.tags...))
	switch a {
	case "true":
		assert.Assert(s, res)
	case "false":
		assert.Assert(s, !res)
	default:
		s.Fatalf("unexpected")
	}
}
