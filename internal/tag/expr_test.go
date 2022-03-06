package tag

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestParseExpr(t *testing.T) {
	//expr, err := Parse("@abc or @feg and @def or (@abc and not @xyz)")
	expr, err := Parse("@def and (@xyz or @bdf) and not @qzy or @abc")
	assert.NilError(t, err)
	t.Logf("%s", expr)
	assert.Assert(t, expr.Match(NewTags("@abc")))
	assert.Assert(t, expr.Match(NewTags("@def", "@xyz")))
	assert.Assert(t, expr.Match(NewTags("@def", "@bdf")))
	assert.Assert(t, !expr.Match(NewTags("@def", "@xyz", "@qzy")))
}
