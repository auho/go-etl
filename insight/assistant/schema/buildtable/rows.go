package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/schema"
)

type RowsTable struct {
	table
	rows assistant.Entity
}

func NewRowsTable(rows assistant.Entity, opts ...TableOption) *RowsTable {
	t := &RowsTable{}
	t.rows = rows
	t.db = rows.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *RowsTable) build() {
	t.initCommand(t.rows.TableName())
	t.AddPKBigInt(t.rows.IDName())

	t.execRawCommandFunc(t.rows)
}

func (t *RowsTable) WithCommand(fn func(*schema.Command)) *RowsTable {
	fn(t.Command)

	return t
}
