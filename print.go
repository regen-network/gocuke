package gocuke

import (
	"fmt"
	"github.com/cucumber/messages-go/v16"
)

func printAstNode(x interface{}) string {
	if x == nil {
		return "nil"
	}

	switch x := x.(type) {
	case *messages.Scenario:
		str := fmt.Sprintf(`%s: %s`, x.Keyword, x.Name)
		return str
	default:
		return fmt.Sprintf("unexpected type %T", x)
	}
}
