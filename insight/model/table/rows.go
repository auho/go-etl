package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type RowsTable struct {
	table
	rows model.Rowsor
}

func NewRowsTable(rows model.Rowsor) *RowsTable {
	t := &RowsTable{}
	t.rows = rows

	t.buildData()

	return t
}

func (t *RowsTable) buildData() {
	t.initCommand(t.rows.TableName())
}
