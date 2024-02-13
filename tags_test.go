package gocuke

import (
	"strings"
	"testing"

	tagexpressions "github.com/cucumber/tag-expressions/go/v6"
	"gotest.tools/v3/assert"

	"github.com/regen-network/gocuke/internal/tag"
)

func TestTags(t *testing.T) {
	NewRunner(t, &tagsSuite{}).Path("features/tags.feature").Run()
}

type tagsSuite struct {
	TestingT
	expr tagexpressions.Evaluatable
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
	s.expr, err = tagexpressions.Parse(a.Content)
	assert.NilError(s, err, a.Content)
}

func (s *tagsSuite) IMatch(a string) {
	s.tags = strings.Split(a, " ")
}

func (s *tagsSuite) TheResultIs(a string) {
	res := tag.NewTags(s.tags...).Match(s.expr)
	switch a {
	case "true":
		assert.Assert(s, res)
	case "false":
		assert.Assert(s, !res)
	default:
		s.Fatalf("unexpected")
	}
}
