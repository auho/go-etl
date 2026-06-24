package altertable

import (
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
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

func (t *Table) WithCommand(fn func(command *schema.Command)) *Table {
	t.commandFunc = fn

	return t
}
