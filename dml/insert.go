package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
)

type Insert struct {
	name         string
	fields       []string
	commander    command.InsertCommander
	sqlCommander Query
}

func NewInsert(name string, q Query) *Insert {
	i := &Insert{}
	i.name = name
	i.sqlCommander = q
	i.commander = newInsertCommand()

	return i
}

func NewInsertWithSelectFields(name string, q Query) *Insert {
	i := NewInsert(name, q)
	i.fields = q.FieldsForInsert()

	return i
}

func NewInsertWithFields(name string, q Query, fields []string) *Insert {
	i := NewInsert(name, q)
	i.fields = fields

	return i
}

func (i *Insert) prepare() {
	i.commander.SetTable(i.name)
	i.commander.SetFields(i.fields)
}

func (i *Insert) Sql() string {
	i.prepare()

	return fmt.Sprintf("%s%s", i.commander.Insert(), i.sqlCommander.Sql())
}
