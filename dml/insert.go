package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
)

type Insert struct {
	name         string
	fields       []string
	commander    command.InsertCommander
	sqlCommander command.SqlCommander
}

func NewInsert(name string, s command.SqlCommander, fields []string) *Insert {
	i := &Insert{}
	i.name = name
	i.fields = fields
	i.sqlCommander = s
	i.commander = newInsertCommand()

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
