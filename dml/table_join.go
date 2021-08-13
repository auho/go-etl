package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
)

type TableJoin struct {
	commander     command.DriverCommander
	tables        []*Table
	limit         []int
	toStringFuncs map[string]func() string
}

func NewTableJoin() *TableJoin {
	tj := &TableJoin{}
	tj.tables = make([]*Table, 0)
	tj.limit = make([]int, 0)
	tj.toStringFuncs = make(map[string]func() string)
	tj.commander = newDriverCommand()

	tj.init()

	return tj
}

func (tj *TableJoin) init() {
	tj.toStringFuncs[reservedSelect] = func() string {
		return tj.commander.SelectToString(tj.mergeTable(tj.tables, func(t *Table) []string {
			return t.commander.BuildSelect()
		}))
	}

	tj.toStringFuncs[reservedFrom] = func() string {
		return tj.commander.FromToString(tj.mergeTable(tj.tables, func(t *Table) []string {
			return t.commander.BuildFrom()
		}))
	}

	tj.toStringFuncs[reservedWhere] = func() string {
		return tj.commander.WhereToString(tj.mergeTable(tj.tables, func(t *Table) []string {
			return t.commander.BuildWhere()
		}))
	}

	tj.toStringFuncs[reservedGroupBy] = func() string {
		return tj.commander.GroupByToString(tj.mergeTable(tj.tables, func(t *Table) []string {
			return t.commander.BuildGroupBy()
		}))
	}

	tj.toStringFuncs[reservedOrderBy] = func() string {
		return tj.commander.OrderByToString(tj.mergeTable(tj.tables, func(t *Table) []string {
			return t.commander.BuildOrderBy()
		}))
	}

	tj.toStringFuncs[reservedLimit] = func() string {
		return tj.commander.LimitToString(tj.limit)
	}
}

func (tj *TableJoin) Table(t *Table) *TableJoin {
	tj.addTable(t)

	return tj
}

func (tj *TableJoin) LeftJoin(t *Table, keys []string, joinTable *Table, joinTableKeys []string) *TableJoin {
	if joinTable == nil {
		joinTable = tj.tables[len(tj.tables)-1]
	}

	if joinTableKeys == nil {
		joinTableKeys = keys
	}

	t.LeftJoin(keys, joinTable, joinTableKeys)
	tj.addTable(t)

	return tj
}

func (tj *TableJoin) Limit(start int, offset int) *TableJoin {
	tj.limit = []int{start, offset}

	return tj
}

func (tj *TableJoin) Sql() string {
	tj.prepare()

	ss := tj.runToStringFuncs([]string{
		reservedSelect,
		reservedFrom,
		reservedWhere,
		reservedGroupBy,
		reservedOrderBy,
		reservedLimit,
	})

	return fmt.Sprintf("%s%s%s%s%s%s", ss...)
}

func (tj *TableJoin) prepare() {
	for _, t := range tj.tables {
		t.Prepare()
	}
}

func (tj *TableJoin) addTable(t *Table) {
	tj.tables = append(tj.tables, t)
}

func (tj *TableJoin) mergeTable(ts []*Table, f func(table *Table) []string) []string {
	s := make([]string, 0)
	for _, t := range ts {
		s = append(s, f(t)...)
	}

	return s
}
func (tj *TableJoin) mergeSlice(ss [][]string) []string {
	s := make([]string, 0)
	for _, v := range ss {
		s = append(s, v...)
	}

	return s
}

func (tj *TableJoin) runToStringFuncs(ns []string) []interface{} {
	ss := make([]interface{}, 0)
	for _, n := range ns {
		ss = append(ss, tj.runToStringFunc(n))
	}

	return ss
}

func (tj *TableJoin) runToStringFunc(n string) string {
	f := tj.toStringFuncs[n]
	return f()
}
