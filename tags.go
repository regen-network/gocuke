package gocuke

import (
	"flag"
	"fmt"
	"github.com/aaronc/gocuke/internal/tag"
)

var flagTags = flag.String("gocuke.tags", "",
	"specify a cucumber tags expression to select tests to run (ex. 'not @long')")

var globalTagExpr *tag.Expr

func initGlobalTagExpr() *tag.Expr {
	if globalTagExpr == nil {
		if flagTags != nil && *flagTags != "" {
			var err error
			globalTagExpr, err = tag.ParseExpr(*flagTags)
			if err != nil {
				if err != nil {
					panic(fmt.Errorf("error parsing tag expression %q: %v\n", *flagTags, err))
				}
			}
		}
	}

	return globalTagExpr
}

// Tags will run only the tests selected by the provided tag expression.
// Tags expressions use the keywords and, or and not and allow expressions
// in parentheses to allow expressions like "(@smoke or @ui) and (not @slow)".
func (r *Runner) Tags(tagExpr string) *Runner {
	var err error
	r.tagExpr, err = tag.ParseExpr(tagExpr)
	if err != nil {
		r.topLevelT.Fatalf("error parsing tag expression %s: %v", tagExpr, err)
	}
	return r
}

// ShortTags specifies which tag expression will be used to select for tests
// when testing.Short() is true. This tag expression will be combined with
// any other tag expression that is applied with Tags() when running short tests.
func (r *Runner) ShortTags(tagExpr string) *Runner {
	var err error
	r.shortTagExpr, err = tag.ParseExpr(tagExpr)
	if err != nil {
		r.topLevelT.Fatalf("error parsing tag expression %s: %v", tagExpr, err)
	}
	return r
}
