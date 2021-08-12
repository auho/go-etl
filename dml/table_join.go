package dml

import "fmt"

type TableJoin struct {
	commander driverCommander
	tables    []*Table
	limit     []int
}

func NewTableJoin() *TableJoin {
	tj := &TableJoin{}
	tj.tables = make([]*Table, 0)
	tj.limit = make([]int, 0)

	tj.commander = newDriverCommand()

	return tj
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

	s := tj.commander.SelectToString(tj.mergeTable(tj.tables, func(t *Table) []string {
		return t.commander.BuildSelect(t.fields)
	}))

	f := tj.commander.FromToString(tj.mergeTable(tj.tables, func(t *Table) []string {
		return t.commander.BuildFrom(t.join)
	}))

	w := tj.commander.WhereToString(tj.mergeTable(tj.tables, func(t *Table) []string {
		return t.commander.BuildWhere(t.where)
	}))

	g := tj.commander.GroupByToString(tj.mergeTable(tj.tables, func(t *Table) []string {
		return t.commander.BuildGroupBy(t.groupBy)
	}))

	o := tj.commander.OrderByToString(tj.mergeTable(tj.tables, func(t *Table) []string {
		return t.commander.BuildOrderBy(t.orderBy)
	}))

	l := tj.commander.LimitToString(tj.limit)

	return fmt.Sprintf("%s%s%s%s%s%s", s, f, w, g, o, l)
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
