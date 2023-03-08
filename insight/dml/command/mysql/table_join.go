package mysql

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/dml/command"
)

var _ command.TableJoinCommander = (*tableJoinCommand)(nil)

type tableJoinCommand struct {
	mysql
	toStringFuncs map[string]func() string
	commands      []command.TableCommander
	limit         []int
}

func NewTableJoinCommand() *tableJoinCommand {
	c := &tableJoinCommand{}
	c.init()

	return c
}

func (c *tableJoinCommand) SetCommands(cs []command.TableCommander) {
	c.commands = cs
}

func (c *tableJoinCommand) SetLimit(l []int) {
	c.limit = l
}

func (c *tableJoinCommand) BuildFieldsForInsert() []string {
	return c.mergeCommand(c.commands, func(tc command.TableCommander) []string {
		return tc.BuildFieldsForInsert()
	})
}

func (c *tableJoinCommand) Query() string {
	ss := c.runToStringFuncs([]string{
		command.ReservedSelect,
		command.ReservedFrom,
		command.ReservedWhere,
		command.ReservedGroupBy,
		command.ReservedOrderBy,
		command.ReservedLimit,
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
			command.ReservedFrom,
			command.ReservedSet,
			command.ReservedWhere,
			command.ReservedOrderBy,
			command.ReservedLimit,
		})...,
	)

	return fmt.Sprintf("UPDATE %s %s%s%s%s%s", ss...)
}

func (c *tableJoinCommand) DeleteQuery() string {
	ss := append([]any{
		c.addBackQuote(c.commands[0].Name())},
		c.runToStringFuncs([]string{
			command.ReservedFrom,
			command.ReservedWhere,
			command.ReservedOrderBy,
			command.ReservedLimit,
		})...,
	)

	return fmt.Sprintf("DELETE %s %s%s%s%s", ss...)
}

func (c *tableJoinCommand) init() {
	c.toStringFuncs = make(map[string]func() string)

	c.toStringFuncs[command.ReservedSelect] = func() string {
		return c.SelectToString(c.mergeCommand(c.commands, func(tc command.TableCommander) []string {
			return tc.BuildSelect()
		}))
	}

	c.toStringFuncs[command.ReservedFrom] = func() string {
		return c.FromToString(c.mergeCommand(c.commands, func(tc command.TableCommander) []string {
			return tc.BuildFrom()
		}))
	}

	c.toStringFuncs[command.ReservedWhere] = func() string {
		return c.WhereToString(c.mergeCommand(c.commands, func(tc command.TableCommander) []string {
			return tc.BuildWhere()
		}))
	}

	c.toStringFuncs[command.ReservedGroupBy] = func() string {
		return c.GroupByToString(c.mergeCommand(c.commands, func(tc command.TableCommander) []string {
			return tc.BuildGroupBy()
		}))
	}

	c.toStringFuncs[command.ReservedOrderBy] = func() string {
		return c.OrderByToString(c.mergeCommand(c.commands, func(tc command.TableCommander) []string {
			return tc.BuildOrderBy()
		}))
	}

	c.toStringFuncs[command.ReservedLimit] = func() string {
		return c.LimitToString(c.limit)
	}

	c.toStringFuncs[command.ReservedSet] = func() string {
		return c.SetToString(c.mergeCommand(c.commands, func(tc command.TableCommander) []string {
			return tc.BuildSet()
		}))
	}
}

func (c *tableJoinCommand) mergeCommand(ts []command.TableCommander, f func(table command.TableCommander) []string) []string {
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
