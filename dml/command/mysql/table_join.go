package mysql

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
)

type tableJoinCommand struct {
	mysql
	toStringFuncs map[string]func() string
	commands      []*tableCommand
	limit         []int
}

func NewTableJoinCommand() *tableJoinCommand {
	c := &tableJoinCommand{}
	c.init()

	return c
}

func (c *tableJoinCommand) SetCommands(cs []*tableCommand) {
	c.commands = cs
}

func (c *tableJoinCommand) SetLimit(l []int) {
	c.limit = l
}

func (c *tableJoinCommand) BuildFieldsForInsert() []string {
	return c.mergeCommand(c.commands, func(tc *tableCommand) []string {
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

func (c *tableJoinCommand) init() {
	c.toStringFuncs = make(map[string]func() string)

	c.toStringFuncs[command.ReservedSelect] = func() string {
		return c.SelectToString(c.mergeCommand(c.commands, func(tc *tableCommand) []string {
			return tc.BuildSelect()
		}))
	}

	c.toStringFuncs[command.ReservedFrom] = func() string {
		return c.FromToString(c.mergeCommand(c.commands, func(tc *tableCommand) []string {
			return tc.BuildFrom()
		}))
	}

	c.toStringFuncs[command.ReservedWhere] = func() string {
		return c.WhereToString(c.mergeCommand(c.commands, func(tc *tableCommand) []string {
			return tc.BuildWhere()
		}))
	}

	c.toStringFuncs[command.ReservedGroupBy] = func() string {
		return c.GroupByToString(c.mergeCommand(c.commands, func(tc *tableCommand) []string {
			return tc.BuildGroupBy()
		}))
	}

	c.toStringFuncs[command.ReservedOrderBy] = func() string {
		return c.OrderByToString(c.mergeCommand(c.commands, func(tc *tableCommand) []string {
			return tc.BuildOrderBy()
		}))
	}

	c.toStringFuncs[command.ReservedLimit] = func() string {
		return c.LimitToString(c.limit)
	}
}

func (c *tableJoinCommand) mergeCommand(ts []*tableCommand, f func(table *tableCommand) []string) []string {
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

func (c *tableJoinCommand) runToStringFuncs(ns []string) []interface{} {
	ss := make([]interface{}, 0)
	for _, n := range ns {
		ss = append(ss, c.runToStringFunc(n))
	}

	return ss
}

func (c *tableJoinCommand) runToStringFunc(n string) string {
	f := c.toStringFuncs[n]
	return f()
}
