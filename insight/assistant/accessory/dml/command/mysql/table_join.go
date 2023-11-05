package mysql

import (
	"fmt"

	command2 "github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
)

var _ command2.TableJoinCommander = (*tableJoinCommand)(nil)

type tableJoinCommand struct {
	mysql
	toStringFuncs map[string]func() string
	commands      []command2.TableCommander
	limit         []int
}

func NewTableJoinCommand() *tableJoinCommand {
	c := &tableJoinCommand{}
	c.init()

	return c
}

func (c *tableJoinCommand) SetCommands(cs []command2.TableCommander) {
	c.commands = cs
}

func (c *tableJoinCommand) SetLimit(l []int) {
	c.limit = l
}

func (c *tableJoinCommand) BuildFieldsForInsert() []string {
	return c.mergeCommand(c.commands, func(tc command2.TableCommander) []string {
		return tc.BuildFieldsForInsert()
	})
}

func (c *tableJoinCommand) Query() string {
	ss := c.runToStringFuncs([]string{
		command2.ReservedSelect,
		command2.ReservedFrom,
		command2.ReservedWhere,
		command2.ReservedGroupBy,
		command2.ReservedOrderBy,
		command2.ReservedLimit,
	})

	return fmt.Sprintf("%s%s%s%s%s%s", ss...)
}

func (c *tableJoinCommand) InsertQuery(name string) string {
	return c.mysql.insert(name, c)
}

func (c *tableJoinCommand) InsertWithFieldsQuery(name string, fields []string) string {
	return c.mysql.insertWithFields(name, fields, c)
}

func (c *tableJoinCommand) UpdateQuery() string {
	ss := append([]any{
		c.addBackQuote(c.commands[0].Name())},
		c.runToStringFuncs([]string{
			command2.ReservedFrom,
			command2.ReservedSet,
			command2.ReservedWhere,
			command2.ReservedOrderBy,
			command2.ReservedLimit,
		})...,
	)

	return fmt.Sprintf("UPDATE %s %s%s%s%s%s", ss...)
}

func (c *tableJoinCommand) DeleteQuery() string {
	ss := append([]any{
		c.addBackQuote(c.commands[0].Name())},
		c.runToStringFuncs([]string{
			command2.ReservedFrom,
			command2.ReservedWhere,
			command2.ReservedOrderBy,
			command2.ReservedLimit,
		})...,
	)

	return fmt.Sprintf("DELETE %s %s%s%s%s", ss...)
}

func (c *tableJoinCommand) init() {
	c.toStringFuncs = make(map[string]func() string)

	c.toStringFuncs[command2.ReservedSelect] = func() string {
		return c.SelectToString(c.mergeCommand(c.commands, func(tc command2.TableCommander) []string {
			return tc.BuildSelect()
		}))
	}

	c.toStringFuncs[command2.ReservedFrom] = func() string {
		return c.FromToString(c.mergeCommand(c.commands, func(tc command2.TableCommander) []string {
			return tc.BuildFrom()
		}))
	}

	c.toStringFuncs[command2.ReservedWhere] = func() string {
		return c.WhereToString(c.mergeCommand(c.commands, func(tc command2.TableCommander) []string {
			return tc.BuildWhere()
		}))
	}

	c.toStringFuncs[command2.ReservedGroupBy] = func() string {
		return c.GroupByToString(c.mergeCommand(c.commands, func(tc command2.TableCommander) []string {
			return tc.BuildGroupBy()
		}))
	}

	c.toStringFuncs[command2.ReservedOrderBy] = func() string {
		return c.OrderByToString(c.mergeCommand(c.commands, func(tc command2.TableCommander) []string {
			return tc.BuildOrderBy()
		}))
	}

	c.toStringFuncs[command2.ReservedLimit] = func() string {
		return c.LimitToString(c.limit)
	}

	c.toStringFuncs[command2.ReservedSet] = func() string {
		return c.SetToString(c.mergeCommand(c.commands, func(tc command2.TableCommander) []string {
			return tc.BuildSet()
		}))
	}
}

func (c *tableJoinCommand) mergeCommand(ts []command2.TableCommander, f func(table command2.TableCommander) []string) []string {
	s := make([]string, 0)
	for _, t := range ts {
		s = append(s, f(t)...)
	}

	return s
}

func (c *tableJoinCommand) mergeSlice(ss [][]string) []string {
	s := make([]string, 0)
	for _, v := range ss {
		s = append(s, v...)
	}

	return s
}

func (c *tableJoinCommand) runToStringFuncs(ns []string) []any {
	ss := make([]any, 0)
	for _, n := range ns {
		ss = append(ss, c.runToStringFunc(n))
	}

	return ss
}

func (c *tableJoinCommand) runToStringFunc(n string) string {
	f := c.toStringFuncs[n]
	return f()
}
