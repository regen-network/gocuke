package gocuke

import (
	"fmt"
	"github.com/cockroachdb/apd/v3"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
	"math/big"
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
	assert.DeepEqual(s, s.orig, s.parsed, decComparer, bigIntComparer)
}

var decComparer = cmp.Comparer(func(x, y *apd.Decimal) bool {
	return x.Cmp(y) == 0
})

var bigIntComparer = cmp.Comparer(func(x, y *big.Int) bool {
	return x.Cmp(y) == 0
})

func (s *valuesSuite) AnyInt64String(t *rapid.T) {
	s.orig = rapid.Int64().Draw(t, "orig")
	s.str = fmt.Sprintf("%d", s.orig)
}

func (s *valuesSuite) WhenIConvertItToAnInt64() {
	s.parsed = toInt64(s, s.str)
}

var decGen = rapid.Custom(func(t *rapid.T) *apd.Decimal {
	nBytes := rapid.IntRange(1, 16).Draw(t, "nBytes").(int)
	bytes := make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		bytes[i] = rapid.Byte().Draw(t, fmt.Sprintf("byte%d", i)).(byte)
	}
	coeff := &apd.BigInt{}
	coeff.SetBytes(bytes)
	neg := rapid.Bool().Draw(t, "neg").(bool)
	if neg {
		coeff = coeff.Neg(coeff)
	}
	exp := rapid.Int32Range(-5000, 5000).Draw(t, "exp").(int32)
	return apd.NewWithBigInt(coeff, exp)
})

func (s *valuesSuite) AnyDecimalString(t *rapid.T) {
	x := decGen.Draw(t, "x").(*apd.Decimal)
	s.orig = x
	s.str = x.String()
}

func (s *valuesSuite) WhenIConvertItToADecimal() {
	s.parsed = toDecimal(s, s.str)
}

var bigIntGen = rapid.Custom(func(t *rapid.T) *big.Int {
	nBytes := rapid.IntRange(1, 16).Draw(t, "nBytes").(int)
	bytes := make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		bytes[i] = rapid.Byte().Draw(t, fmt.Sprintf("byte%d", i)).(byte)
	}
	x := &big.Int{}
	x.SetBytes(bytes)
	neg := rapid.Bool().Draw(t, "neg").(bool)
	if neg {
		x = x.Neg(x)
	}
	return x
})

func (s *valuesSuite) AnyBigIntegerString(t *rapid.T) {
	x := bigIntGen.Draw(t, "x").(*big.Int)
	s.orig = x
	s.str = x.String()
}

func (s *valuesSuite) WhenIConvertItToABigInteger() {
	s.parsed = toBigInt(s, s.str)
}