package datatable

import (
	"testing"

	"github.com/regen-network/gocuke"
	"gotest.tools/v3/assert"
)

func TestDataTable(t *testing.T) {
	gocuke.NewRunner(t, &dataTableSuite{}).
		Path("datatable.feature").Run()
}

type dataTableSuite struct {
	gocuke.TestingT
	datatable gocuke.DataTable
	total     int64
}

func (d *dataTableSuite) Before(t gocuke.TestingT) {
	d.TestingT = t
}

func (s *dataTableSuite) ThisDataTable(a gocuke.DataTable) {
	s.datatable = a
}

func (s *dataTableSuite) ItHasRowsAndColumns(rows int64, cols int64) {
	assert.Equal(s, int(rows), s.datatable.NumRows())
	assert.Equal(s, int(cols), s.datatable.NumCols())
}

func (s *dataTableSuite) RowsAsAHeaderTable(rows int64) {
	assert.Equal(s, int(rows), s.datatable.HeaderTable().NumRows())
}

func (s *dataTableSuite) TheValuesAddUpWhenAccessedAsAHeaderTable() {
	ht := s.datatable.HeaderTable()
	s.total = 0
	for i := 0; i < ht.NumRows(); i++ {
		sum := ht.Get(i, "x").Int64() + ht.Get(i, "y").Int64() + ht.Get(i, "z").Int64()
		assert.Equal(s, ht.Get(i, "x + y + z").Int64(), sum)
		s.total += sum
	}
}

func (s *dataTableSuite) TheTotalSumOfTheXYZColumnIs(a int64) {
	assert.Equal(s, a, s.total)
}
