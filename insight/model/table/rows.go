package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type RowsTable struct {
	table
	rows model.Rowsor
}

func NewRowsTable(rows model.Rowsor) *RowsTable {
	rt := &RowsTable{}
	rt.rows = rows

	rt.buildData()

	return rt
}

func (dt *RowsTable) buildData() {
	dt.initTable(dt.rows.TableName())
}
