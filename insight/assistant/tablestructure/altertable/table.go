package altertable

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

type Table struct {
	baseTable
}

func NewTable(tableName string) *Table {
	t := &Table{}
	t.baseTable = newBaseTable(tableName)

	return t
}

func (t *Table) Build(db *simpledb.SimpleDB) error {
	return t.build(t.SQL(), db)
}

func (t *Table) BuildChange(db *simpledb.SimpleDB) error {
	return t.build(t.SqlForChange(), db)
}

func (t *Table) WithCommand(fn func(command *tablestructure.Command)) *Table {
	t.commandFunc = fn

	return t
}
