package gocuke

import (
	"github.com/cockroachdb/apd/v3"
	"github.com/cucumber/messages-go/v16"
	"math/big"
	"reflect"
	"testing"
)

type DataTable struct {
	t     *testing.T
	table *messages.PickleTable
}

func (d DataTable) NumRows() int {
	return len(d.table.Rows)
}

func (d DataTable) NumCols() int {
	if len(d.table.Rows) == 0 {
		d.t.Fatalf("no table rows")
	}

	return len(d.table.Rows[0].Cells)
}

func (d DataTable) Cell(row, col int) *Cell {
	if row >= len(d.table.Rows) {
		d.t.Fatalf("table row %d out of range", row)
	}

	r := d.table.Rows[row]

	if col >= len(r.Cells) {
		d.t.Fatalf("table column %d out of range", col)
	}

	return &Cell{
		t:     d.t,
		value: r.Cells[col].Value,
	}
}

type Cell struct {
	t     *testing.T
	value string
}

func (c Cell) String() string {
	return c.value
}

func (c Cell) Int64() int64 {
	return toInt64(c.t, c.value)
}

func (c Cell) Float64() float64 {
	return toFloat64(c.t, c.value)
}

func (c Cell) BigInt() *big.Int {
	return toBigInt(c.t, c.value)
}

func (c Cell) Decimal() *apd.Decimal {
	return toDecimal(c.t, c.value)
}

var dataTableType = reflect.TypeOf(DataTable{})
