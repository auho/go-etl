package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
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

func (t *RowsTable) Exec(fn func(command *tablestructure.Command)) *RowsTable {
	fn(t.Command)

	return t
}
