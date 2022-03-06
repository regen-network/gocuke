package tag

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type Expr = Or

type Or struct {
	LHS *And   `@@`
	RHS []*And `("or" @@)*`
}

type Term struct {
	Tag   string `@Tag |`
	Not   *Term  `"not" @@ |`
	Group *Expr  `"(" @@ ")"`
}

type And struct {
	LHS *Term   `@@`
	RHS []*Term `("and" @@)*`
}

func (s Or) String() string {
	strs := make([]string, 1+len(s.RHS))
	strs[0] = s.LHS.String()
	for i, rh := range s.RHS {
		strs[i+1] = rh.String()
	}
	return strings.Join(strs, " or ")
}

func (a And) String() string {
	strs := make([]string, 1+len(a.RHS))
	strs[0] = a.LHS.String()
	for i, rh := range a.RHS {
		strs[i+1] = rh.String()
	}
	return strings.Join(strs, " and ")
}

func (t Term) String() string {
	if t.Not != nil {
		return "not " + t.Not.String()
	} else if t.Group != nil {
		return "(" + t.Group.String() + ")"
	} else {
		return t.Tag
	}
}

type Tags map[string]bool

func NewTags(tags ...string) Tags {
	res := map[string]bool{}
	for _, tag := range tags {
		res[tag] = true
	}
	return res
}

func (s Or) Match(tags Tags) bool {
	if s.LHS.Match(tags) {
		return true
	}
	for _, rh := range s.RHS {
		if rh.Match(tags) {
			return true
		}
	}
	return false
}

func (a And) Match(tags Tags) bool {
	if !a.LHS.Match(tags) {
		return false
	}
	for _, rh := range a.RHS {
		if !rh.Match(tags) {
			return false
		}
	}
	return true
}

func (t Term) Match(tags Tags) bool {
	if t.Not != nil {
		return !t.Not.Match(tags)
	} else if t.Group != nil {
		return t.Group.Match(tags)
	} else {
		return tags[t.Tag]
	}
}

func Parse(expr string) (*Expr, error) {
	ast := &Expr{}
	err := parser.ParseString("", expr, ast)
	return ast, err
}

var (
	parser = participle.MustBuild(&Expr{},
		participle.Lexer(lex),
		participle.Elide("Whitespace"),
	)
	lex = lexer.MustSimple([]lexer.Rule{
		{"Tag", `@[A-Za-z]\w*`, nil},
		{"Keyword", `and|or|not`, nil},
		{"Paren", `[\(\)]`, nil},
		{"Whitespace", `[ \t\n\r]+`, nil},
	})
)
