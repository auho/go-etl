package dml

import (
	"github.com/auho/go-etl/dml/command"
)

type Insert struct {
	name      string
	fields    []string
	commander command.InsertCommander
	q         command.Query
}

func NewInsert(name string, q command.Query) *Insert {
	i := &Insert{}
	i.name = name
	i.q = q
	i.commander = newInsertCommand()

	return i
}

func NewInsertWithFields(name string, q command.Query, fields []string) *Insert {
	i := NewInsert(name, q)
	i.fields = fields

	return i
}

func NewInsertWithSelectFields(name string, q command.Query) *Insert {
	i := NewInsert(name, q)

	return i
}

func (i *Insert) Sql() string {
	i.prepare()

	return i.commander.Insert()
}

func (i *Insert) SqlWithFields() string {
	i.prepare()

	return i.commander.InsertWithFields()
}

func (i *Insert) prepare() {
	i.commander.SetTable(i.name)
	i.commander.SetFields(i.fields)
	i.commander.SetQuery(i.q)
}
