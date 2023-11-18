package dml

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
)

var _ Tabler = (*TableJoin)(nil)

type TableJoin struct {
	commander command.TableJoinCommander
	tables    []*Table
	limit     []int
	set       []*command.Set
}

func newTableJoin(driver string) *TableJoin {
	tj := &TableJoin{}
	tj.commander = newTableJoinCommand(driver)
	tj.tables = make([]*Table, 0)
	tj.limit = make([]int, 0)
	tj.set = make([]*command.Set, 0)

	return tj
}

func NewTableJoin() *TableJoin {
	return newTableJoin("")
}

func (tj *TableJoin) Table(t *Table) *TableJoin {
	tj.addTable(t)

	return tj
}

// LeftJoin
//
// left table 为 nil，默认取上一个设定的 table, table fields
//
// t => right table
// fields => right table join fields
// leftTable => left table
// leftFields => left table join fields
func (tj *TableJoin) LeftJoin(rightTable *Table, rightFields []string, leftTable *Table, leftFields []string) *TableJoin {
	if leftTable == nil {
		leftTable = tj.tables[len(tj.tables)-1]
	}

	if leftFields == nil {
		leftFields = rightFields
	}

	if len(rightFields) <= 0 || len(leftFields) <= 0 {
		panic("left join fields not found")
	}

	rightTable.LeftJoin(rightFields, leftTable, leftFields)
	tj.addTable(rightTable)

	return tj
}

func (tj *TableJoin) SetField(t *Table, fields []string, setTable *Table, setFields []string) *TableJoin {
	if setTable == nil {
		setTable = tj.tables[len(tj.tables)-1]
	}

	if setFields == nil {
		setFields = fields
	}

	t.SetSet(command.NewSetField(t.name, fields, setTable.name, setFields))

	return tj
}

func (tj *TableJoin) SetExpression(t *Table, fields []string, setTable *Table, expression []string) *TableJoin {
	if setTable == nil {
		setTable = tj.tables[len(tj.tables)-1]
	}

	if expression == nil {
		panic("set expression not found")
	}

	t.SetSet(command.NewSetExpression(t.name, fields, setTable.name, expression))

	return tj
}

func (tj *TableJoin) SetValue(t *Table, fields []string, setTable *Table, values []any) *TableJoin {
	if setTable == nil {
		setTable = tj.tables[len(tj.tables)-1]
	}

	if values == nil {
		panic("set values not found")
	}

	t.SetSet(command.NewSetValue(t.name, fields, setTable.name, values))

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

func (tj *TableJoin) GetSelectFields() []string {
	var fields []string
	for _, table := range tj.tables {
		fields = append(fields, table.GetSelectFields()...)
	}

	return fields
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
