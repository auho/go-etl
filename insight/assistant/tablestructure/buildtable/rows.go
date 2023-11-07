package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
)

type RowsTable struct {
	table
	rows assistant.Rowsor
}

func NewRowsTable(rows assistant.Rowsor, opts ...TableOption) *RowsTable {
	t := &RowsTable{}
	t.rows = rows
	t.db = rows.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *RowsTable) build() {
	t.initCommand(t.rows.TableName())
	t.AddPkBigInt(t.rows.GetIdName())

	t.execCommand()
	t.execRowsCommand(t.rows)
}
