package dml

import (
	"github.com/auho/go-etl/v2/insight/dml/command"
)

type TableJoin struct {
	commander command.TableJoinCommander
	tables    []*Table
	limit     []int
	set       []*command.Set
}

func NewTableJoin() *TableJoin {
	tj := &TableJoin{}
	tj.commander = newTableJoinCommand()
	tj.tables = make([]*Table, 0)
	tj.limit = make([]int, 0)
	tj.set = make([]*command.Set, 0)

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

func (tj *TableJoin) Set(t *Table, keys []string, setTable *Table, setKeys []string) *TableJoin {
	if setTable == nil {
		setTable = tj.tables[len(tj.tables)-1]
	}

	if setKeys == nil {
		setKeys = keys
	}

	t.SetSet(command.NewSet(t.name, keys, setTable.name, setKeys))

	return tj
}

func (tj *TableJoin) SetExpression(t *Table, keys []string, setTable *Table, setKeys []string) *TableJoin {
	if setTable == nil {
		setTable = tj.tables[len(tj.tables)-1]
	}

	if setKeys == nil {
		setKeys = keys
	}

	t.SetSet(command.NewExpressionSet(t.name, keys, setTable.name, setKeys))

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

func (tj *TableJoin) InsertSql(name string) string {
	tj.prepare()

	return tj.commander.InsertQuery(name)
}

func (tj *TableJoin) InsertWithFieldsSql(name string, fields []string) string {
	tj.prepare()

	return tj.commander.InsertWithFieldsQuery(name, fields)
}

func (tj *TableJoin) UpdateSql() string {
	tj.prepare()

	return tj.commander.UpdateQuery()
}

func (tj *TableJoin) DeleteSql() string {
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
	tj.commander.SetLimit(tj.limit)
}

func (tj *TableJoin) addTable(t *Table) {
	tj.tables = append(tj.tables, t)
}
