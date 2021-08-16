package mysql

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/dml/command"
)

type insertCommand struct {
	mysql
	name          string
	nameBackQuote string
	fields        []string
	q             command.Query
}

func NewInsertCommand() *insertCommand {
	return &insertCommand{}
}

func (i *insertCommand) SetTable(name string) {
	i.name = name
	i.nameBackQuote = i.addBackQuote(i.name)
}

func (i *insertCommand) SetFields(fields []string) {
	i.fields = fields
}

func (i *insertCommand) SetQuery(q command.Query) {
	i.q = q
}

func (i *insertCommand) Insert() string {
	return fmt.Sprintf("INSERT INTO %s %s", i.nameBackQuote, i.q.Sql())
}

func (i *insertCommand) InsertWithFields() string {
	s := ""
	fields := make([]string, 0)
	if i.fields == nil {
		i.fields = i.q.FieldsForInsert()
	}

	for _, f := range i.fields {
		fields = append(fields, i.addBackQuote(f))
	}

	s = fmt.Sprintf(" (%s) ", strings.Join(fields, ", "))

	return fmt.Sprintf("INSERT INTO %s%s %s", i.nameBackQuote, s, i.q.Sql())
}
