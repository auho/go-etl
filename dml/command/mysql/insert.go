package mysql

import (
	"fmt"
	"strings"
)

type insertCommand struct {
	mysql
	name          string
	nameBackQuote string
	fields        []string
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

func (i *insertCommand) Insert() string {
	s := ""
	fields := make([]string, 0)
	if i.fields != nil {
		for _, f := range i.fields {
			fields = append(fields, i.addBackQuote(f))
		}

		s = fmt.Sprintf(" (%s) ", strings.Join(fields, ", "))
	}

	return fmt.Sprintf("INSERT INTO %s%s ", i.nameBackQuote, s)
}
