package tag

import (
	messages "github.com/cucumber/messages/go/v22"
	tag "github.com/cucumber/tag-expressions/go/v6"
)

func NewTags(tags ...string) Tags {
	have := map[string]bool{}
	var res []string
	for _, t := range tags {
		if !have[t] {
			have[t] = true
			res = append(res, t)
		}
	}
	return res
}

func NewTagsFromPickleTags(pickleTags []*messages.PickleTag) Tags {
	have := map[string]bool{}
	var res []string
	for _, t := range pickleTags {
		if !have[t.Name] {
			have[t.Name] = true
			res = append(res, t.Name)
		}
	}
	return res
}

type Tags []string

func (t Tags) Match(expr tag.Evaluatable) bool {
	return expr.Evaluate(t)
}
