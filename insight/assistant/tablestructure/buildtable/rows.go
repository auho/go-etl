package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
)

type RowsTable struct {
	table
	rows assistant.Rowsor
}

func NewRowsTable(rows assistant.Rowsor) *RowsTable {
	t := &RowsTable{}
	t.rows = rows

	t.build()

	return t
}

func (t *RowsTable) build() {
	t.initCommand(t.rows.TableName())
	t.AddPkBigInt(t.rows.GetIdName())

	t.execRowsCommand(t.rows)
}
