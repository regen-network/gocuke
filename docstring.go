package gocuke

import "reflect"

type DocString struct {
	MediaType string
	Content   string
}

var docStringType = reflect.TypeOf(DocString{})
