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
	NewRunner(t, &valuesSuite{}).
		Path("features/values.feature").
		ShortTags("not @long").
		Run()
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
	s.AnInt64(rapid.Int64().Draw(t, "orig").(int64))
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
	s.ADecimal(decGen.Draw(t, "x").(*apd.Decimal))
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
	s.ABigInteger(bigIntGen.Draw(t, "x").(*big.Int))
}

func (s *valuesSuite) WhenIConvertItToABigInteger() {
	s.parsed = toBigInt(s, s.str)
}

func (s *valuesSuite) AnInt64(a int64) {
	s.orig = a
	s.str = fmt.Sprintf("%d", a)
}

func (s *valuesSuite) ADecimal(a *apd.Decimal) {
	s.orig = a
	s.str = a.String()
}

func (s *valuesSuite) ABigInteger(a *big.Int) {
	s.orig = a
	s.str = a.String()
}
