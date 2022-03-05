package gocuke

import "testing"

func TestGuessMethodSig(t *testing.T) {
	t.Log(guessMethodSig("I have one"))
	t.Log(guessMethodSig("I have 5"))
	t.Log(guessMethodSig(`I have 5 "foo" and 3 "bar"`))
}
