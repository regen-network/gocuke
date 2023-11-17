package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/regen-network/gocuke"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type suite struct {
	// special arguments like TestingT are injected automatically into exported fields
	gocuke.TestingT
	resp *httptest.ResponseRecorder
}

func TestApi(t *testing.T) {
	scope := &suite{TestingT: t, resp: httptest.NewRecorder()}
	gocuke.NewRunner(t, scope).
		Step(`^I send "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, scope.ISendRequestTo).
		Step(`^the response code should be (\d+)$`, scope.TheResponseCodeShouldBe).
		Step(`^the response should match json:$`, scope.TheResponseShouldMatchJson).
		Run()
}

func (s *suite) ISendRequestTo(method string, endpoint string) {
	req, err := http.NewRequest(method, endpoint, nil)
	assert.Nil(s, err)

	defer func() {
		switch t := recover().(type) {
		case string:
			err = fmt.Errorf(t)
		case error:
			err = t
		}
	}()

	switch endpoint {
	case "/version":
		getVersion(s.resp, req)
	default:
		err = fmt.Errorf("unknown endpoint: %s", endpoint)
	}
	assert.Nil(s, err)
}

func (s *suite) TheResponseCodeShouldBe(code int64) {
	assert.Equalf(s, code, int64(s.resp.Code), "expected response code to be: %d, but actual is: %d", code, s.resp.Code)
}

func (s *suite) TheResponseShouldMatchJson(body gocuke.DocString) {
	require.JSONEq(s, body.Content, s.resp.Body.String())
}
