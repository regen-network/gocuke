package gocuke

import (
	"github.com/cockroachdb/apd/v3"
	"math/big"
	"reflect"
	"strconv"
	"testing"
)

func convertParamValue(t *testing.T, match string, typ reflect.Type) reflect.Value {
	switch typ.Kind() {
	case reflect.Int64:
		return reflect.ValueOf(toInt64(t, match))
	case reflect.Float64:
		return reflect.ValueOf(toFloat64(t, match))
	case reflect.String:
		return reflect.ValueOf(match)
	default:
		if typ == bigIntType {
			return reflect.ValueOf(toBigInt(t, match))
		} else if typ == decType {
			return reflect.ValueOf(toDecimal(t, match))
		} else {
			t.Fatalf("unexpected parameter type %v", typ)
			return reflect.Value{}
		}
	}
}

var (
	bigIntType = reflect.TypeOf(&big.Int{})
	decType    = reflect.TypeOf(&apd.Decimal{})
)

func toInt64(t *testing.T, value string) int64 {
	x, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		t.Fatalf("error converting %s to int64: %v", value, err)
	}
	return x
}

func toFloat64(t *testing.T, value string) float64 {
	x, err := strconv.ParseFloat(value, 64)
	if err != nil {
		t.Fatalf("error converting %s to int: %v", value, err)
	}
	return x
}

func toBigInt(t *testing.T, value string) *big.Int {
	x := &big.Int{}
	x, ok := x.SetString(value, 10)
	if ok {
		t.Fatalf("error converting %s to *big.Int", value)
	}
	return x
}

func toDecimal(t *testing.T, value string) *apd.Decimal {
	x, _, err := apd.NewFromString(value)
	if err != nil {
		t.Fatalf("error converting %s to a decimal: %v", value, err)
	}
	return x
}
