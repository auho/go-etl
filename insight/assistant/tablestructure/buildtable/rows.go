package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
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

	t.execCommandFunc()
	t.execRawCommandFunc(t.rows)
}

func (t *RowsTable) WithCommand(fn func(*tablestructure.Command)) *RowsTable {
	fn(t.Command)

	return t
}
