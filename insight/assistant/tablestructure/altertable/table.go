package altertable

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Table struct {
	baseTable
}

func NewTable(tableName string) *Table {
	t := &Table{}
	t.baseTable = newBaseTable(tableName)

	return t
}

func (t *Table) Build(db *simpleDb.SimpleDB) error {
	return t.build(t.Sql(), db)
}

func (t *Table) BuildChange(db *simpleDb.SimpleDB) error {
	return t.build(t.SqlForChange(), db)
}

func (t *Table) WithCommand(fn func(command *tablestructure.Command)) *Table {
	t.commandFun = fn

	return t
}
