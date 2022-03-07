package gocuke

import "reflect"

// DocString represents a doc string step argument.
type DocString struct {
	MediaType string
	Content   string
}

var docStringType = reflect.TypeOf(DocString{})
