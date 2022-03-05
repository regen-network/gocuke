package gocuke

import (
	"github.com/cucumber/messages-go/v16"
	"testing"
)

func TestGuessMethodSig(t *testing.T) {
	t.Log(guessMethodSigWithText("I have one"))
	t.Log(guessMethodSigWithText("I have 5"))
	t.Log(guessMethodSigWithText(`I have 5 "foo" and 3 "bar"`))
}

func guessMethodSigWithText(text string) methodSig {
	return guessMethodSig(&messages.PickleStep{Text: text})
}
