package gocuke

import (
	"math/big"
	"reflect"

	"github.com/cockroachdb/apd/v3"
	messages "github.com/cucumber/messages/go/v22"
)

// DataTable wraps a data table step argument
type DataTable struct {
	t     TestingT
	table *messages.PickleTable
}

// NumRows returns the number of rows in the data table.
func (d DataTable) NumRows() int {
	return len(d.table.Rows)
}

// NumCols returns the number of columns in the data table.
func (d DataTable) NumCols() int {
	if len(d.table.Rows) == 0 {
		d.t.Fatalf("no table rows")
	}

	return len(d.table.Rows[0].Cells)
}

// Cell returns the cell at the provided 0-based row and col offset.
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

// Cell represents a data table cell.
type Cell struct {
	t     TestingT
	value string
}

// String returns the cell value as a string.
func (c Cell) String() string {
	return c.value
}

// Int64 returns the cell as an int64.
func (c Cell) Int64() int64 {
	return toInt64(c.t, c.value)
}

// BigInt returns the cell as a *big.Int.
func (c Cell) BigInt() *big.Int {
	return toBigInt(c.t, c.value)
}

// Decimal returns the cell value as an *apd.Decimal.
func (c Cell) Decimal() *apd.Decimal {
	return toDecimal(c.t, c.value)
}

// HeaderTable returns the data table as a header table which is a wrapper
// around the table which assumes that the first row is the table header.
func (d DataTable) HeaderTable() *HeaderTable {
	headers := map[string]int{}
	for i := 0; i < d.NumCols(); i++ {
		headers[d.Cell(0, i).String()] = i
	}

	return &HeaderTable{headers: headers, DataTable: d}
}

// HeaderTable is a wrapper around a table which assumes that the first row is \
// the table header.
type HeaderTable struct {
	DataTable
	headers map[string]int
}

// Get returns the cell at the provided row offset (skipping the header row)
// and column name (as indicated in the header). If the column name is not
// found, nil is returned.
func (h *HeaderTable) Get(row int, col string) *Cell {
	if v, ok := h.headers[col]; !ok {
		return nil
	} else {
		return h.DataTable.Cell(row+1, v)
	}
}

// NumRows returns the number of rows in the table (excluding the header row).
func (h *HeaderTable) NumRows() int {
	return h.DataTable.NumRows() - 1
}

var dataTableType = reflect.TypeOf(DataTable{})
