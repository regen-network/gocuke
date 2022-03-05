package gocuke

import (
	"github.com/cockroachdb/apd/v3"
	"github.com/cucumber/messages-go/v16"
	"math/big"
	"reflect"
)

type DataTable struct {
	t     TestingT
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
	t     TestingT
	value string
}

func (c Cell) String() string {
	return c.value
}

func (c Cell) Int64() int64 {
	return toInt64(c.t, c.value)
}

func (c Cell) BigInt() *big.Int {
	return toBigInt(c.t, c.value)
}

func (c Cell) Decimal() *apd.Decimal {
	return toDecimal(c.t, c.value)
}

func (d DataTable) HeaderTable() *HeaderTable {
	headers := map[string]int{}
	for i := 0; i < d.NumCols(); i++ {
		headers[d.Cell(0, i).String()] = i
	}

	return &HeaderTable{headers: headers, DataTable: d}
}

type HeaderTable struct {
	DataTable
	headers map[string]int
}

func (h *HeaderTable) Get(row int, col string) *Cell {
	return h.DataTable.Cell(row+1, h.headers[col])
}

func (h *HeaderTable) NumRows() int {
	return h.DataTable.NumRows() - 1
}

var dataTableType = reflect.TypeOf(DataTable{})
