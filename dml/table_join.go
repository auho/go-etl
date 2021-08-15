package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
)

type TableJoin struct {
	commander     command.TableJoinCommander
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

func (tj *TableJoin) FieldsForInsert() []string {
	return tj.mergeTable(tj.tables, func(t *Table) []string {
		return t.FieldsForInsert()
	})
}

func (tj *TableJoin) Delete() string {
	tj.prepare()

	ss := tj.runToStringFuncs([]string{
		reservedFrom,
		reservedWhere,
		reservedLimit,
	})

	return fmt.Sprintf("%s%s%s", ss...)
}

func (tj *TableJoin) Sql() string {
	tj.prepare()

	return tj.commander.Query()
}

func (tj *TableJoin) prepare() {
	for _, t := range tj.tables {
		t.Prepare()
	}
}

func (tj *TableJoin) addTable(t *Table) {
	tj.tables = append(tj.tables, t)
}
