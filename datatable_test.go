package gocuke

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestDataTable(t *testing.T) {
	NewRunner(t, func(t TestingT) Suite {
		return &dataTableSuite{TestingT: t}
	}).Path("features/datatable.feature").Run()
}

type dataTableSuite struct {
	TestingT
	datatable DataTable
	total     int64
}

func (s *dataTableSuite) ThisDataTable(a DataTable) {
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
