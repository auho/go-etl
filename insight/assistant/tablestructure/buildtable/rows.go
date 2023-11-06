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

	t.buildData()

	return t
}

func (t *RowsTable) buildData() {
	t.initCommand(t.rows.TableName())
}

func (t *RowsTable) Exec(fn func(command *tablestructure.Command)) *RowsTable {
	fn(t.Command)

	return t
}
