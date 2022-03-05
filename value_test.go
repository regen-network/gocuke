package gocuke

import (
	"fmt"
	"github.com/cockroachdb/apd/v3"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
	"pgregory.net/rapid"
	"testing"
)

func TestValues(t *testing.T) {
	NewRunner(t, func(t TestingT) Suite {
		return &valuesSuite{TestingT: t}
	}).WithPath("features/values.feature").Run()
}

type valuesSuite struct {
	TestingT
	orig   interface{}
	str    string
	parsed interface{}
}

func (s *valuesSuite) IGetBackTheOriginalValue() {
	assert.DeepEqual(s, s.orig, s.parsed, decComparer)
}

var decComparer = cmp.Comparer(func(x, y *apd.Decimal) bool {
	return x.Cmp(y) == 0
})

func (s *valuesSuite) AnyInt64String(t *rapid.T) {
	s.orig = rapid.Int64().Draw(t, "orig")
	s.str = fmt.Sprintf("%d", s.orig)
}

func (s *valuesSuite) WhenIConvertItToAnInt64() {
	s.parsed = toInt64(s, s.str)
}

var bigIntGen = rapid.Custom(func(t *rapid.T) *apd.BigInt {
	nBytes := rapid.IntRange(1, 12).Draw(t, "nBytes").(int)
	bytes := make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		bytes[i] = rapid.Byte().Draw(t, fmt.Sprintf("byte%d", i)).(byte)
	}
	x := &apd.BigInt{}
	x.SetBytes(bytes)
	neg := rapid.Bool().Draw(t, "neg").(bool)
	if neg {
		x = x.Neg(x)
	}
	return x
})

func (s *valuesSuite) AnyDecimalString(t *rapid.T) {
	coeff := bigIntGen.Draw(t, "coeff").(*apd.BigInt)
	exp := rapid.Int32Range(apd.MinExponent, apd.MaxExponent).Draw(t, "exp").(int32)
	dec := apd.NewWithBigInt(coeff, exp)
	s.orig = dec
	s.str = dec.String()
}

func (s *valuesSuite) WhenIConvertItToADecimal() {
	s.parsed = toDecimal(s, s.str)
}
