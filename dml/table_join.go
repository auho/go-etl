package dml

import (
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

func (tj *TableJoin) Sql() string {
	tj.prepare()

	return tj.commander.Query()
}

func (tj *TableJoin) Insert(name string) string {
	tj.prepare()

	return tj.commander.InsertQuery(name)
}

func (tj *TableJoin) InsertWithFields(name string, fields []string) string {
	tj.prepare()

	return tj.commander.InsertWithFieldsQuery(name, fields)
}

func (tj *TableJoin) Delete() string {
	tj.prepare()

	return tj.commander.DeleteQuery()
}

func (tj *TableJoin) prepare() {
	commands := make([]command.TableCommander, 0)
	for _, t := range tj.tables {
		t.prepare()

		commands = append(commands, t.commander)
	}

	tj.commander.SetCommands(commands)
}

func (tj *TableJoin) addTable(t *Table) {
	tj.tables = append(tj.tables, t)
}
