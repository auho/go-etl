package dml

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
)

var _ Tabler = (*Table)(nil)

type Table struct {
	commander command.TableCommander
	name      string
	fields    *command.Entities
	where     string
	groupBy   *command.Entities
	orderBy   *command.Entities
	limit     []int
	having    string
	join      *command.Join
	set       []*command.Set
	asSql     string
}

func newTable(name, driver string) *Table {
	t := &Table{}

	t.init(name, driver)

	return t
}

// NewTable
// new sql from table. table as data source
func NewTable(name string) *Table {
	return newTable(name, "")
}

// NewSqlTable
// new sql from sql. Sql result set as data source
// name as the table alias of the result set
func NewSqlTable(name, sql string) *Table {
	t := NewTable(name)
	t.asSql = sql

	t.init(name, "")

	return t
}

func (t *Table) init(name string, driver string) {
	t.name = name
	t.fields = command.NewEntries()
	t.groupBy = command.NewEntries()
	t.orderBy = command.NewEntries()
	t.limit = make([]int, 0)
	t.set = make([]*command.Set, 0)

	t.commander = newTableCommand(driver)
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) Select(fields []string) *Table {
	for _, field := range fields {
		t.fields.AddEntry(field, field)
	}

	return t
}

func (t *Table) SelectAlias(alias map[string]string) *Table {
	for k, v := range alias {
		t.fields.AddEntryExpression(k, v)
	}

	return t
}

func (t *Table) Where(s string) *Table {
	t.where = s

	return t
}

func (t *Table) GroupBy(g []string) *Table {
	for _, v := range g {
		t.groupBy.AddEntry(v, v)
		t.fields.AddEntry(v, v)
	}

	return t
}

func (t *Table) GroupByAlias(g map[string]string) *Table {
	for k, v := range g {
		t.groupBy.AddEntry(k, v)
		t.fields.AddEntry(k, v)
	}

	return t
}

func (t *Table) addOrderBy(flag string, v ...string) *Table {
	v = v[0 : len(v)/2*2]
	for i := 0; i < len(v); i += 2 {
		if flag == command.FlagExpression {
			t.orderBy.AddEntryExpression(v[i], v[i+1])
		} else {
			t.orderBy.AddEntry(v[i], v[i+1])
		}
	}

	return t
}

func (t *Table) OrderBy(v ...string) *Table {
	t.addOrderBy(command.FlagField, v...)

	return t
}

func (t *Table) OrderByExpression(v ...string) *Table {
	t.addOrderBy(command.FlagExpression, v...)

	return t
}

func (t *Table) OrderByAsc(k string) *Table {
	t.OrderBy(k, command.SortAsc)

	return t
}

func (t *Table) OrderByAscExpression(k string) *Table {
	t.OrderByExpression(k, command.SortAsc)

	return t
}

func (t *Table) OrderByDesc(k string) *Table {
	t.OrderBy(k, command.SortDesc)

	return t
}

func (t *Table) OrderByDescExpression(k string) *Table {
	t.OrderByExpression(k, command.SortDesc)

	return t
}

func (t *Table) Having(s string) *Table {
	t.having = s

	return t
}

func (t *Table) Limit(start int, offset int) *Table {
	t.limit = []int{start, offset}

	return t
}

func (t *Table) Aggregation(a map[string]string) *Table {
	for k, v := range a {
		e := command.NewExpressionEntity(k, v)
		t.fields.Add(e)
	}

	return t
}

func (t *Table) addLeftJoin(fields []string, rightTable *Table, rightFields []string) *Table {
	t.join = command.NewLeftJoin(t.GetName(), fields, rightTable.GetName(), rightFields)

	return t
}

// SetField
// update statement set syntax
// map[string]string => map[left table fields]right table fields
func (t *Table) SetField(s map[string]string) *Table {
	fields := make([]string, 0)
	values := make([]string, 0)

	for k, v := range s {
		fields = append(fields, k)
		values = append(values, v)
	}

	t.set = append(t.set, command.NewSetField(t.name, fields, t.name, values))

	return t
}

func (t *Table) SetExpression(s map[string]string) *Table {
	fields := make([]string, 0)
	values := make([]string, 0)

	for k, v := range s {
		fields = append(fields, k)
		values = append(values, v)
	}

	t.set = append(t.set, command.NewSetExpression(t.name, fields, t.name, values))

	return t
}

func (t *Table) SetValue(s map[string]any) *Table {
	fields := make([]string, 0)
	values := make([]any, 0)

	for k, v := range s {
		fields = append(fields, k)
		values = append(values, v)
	}

	t.set = append(t.set, command.NewSetValue(t.name, fields, t.name, values))

	return t
}

func (t *Table) SetSet(s *command.Set) {
	t.set = append(t.set, s)
}

func (t *Table) Sql() string {
	t.prepare()

	return t.commander.Query()
}

func (t *Table) InsertSql(name string) string {
	t.prepare()

	return t.commander.InsertQuery(name)
}

func (t *Table) InsertWithFieldsSql(name string, fields []string) string {
	t.prepare()

	return t.commander.InsertWithFieldsQuery(name, fields)
}

func (t *Table) UpdateSql() string {
	t.prepare()

	return t.commander.UpdateQuery()
}

func (t *Table) DeleteSql() string {
	t.prepare()

	return t.commander.DeleteQuery()
}

func (t *Table) CreateJoin() *TableJoin {
	return newTableJoin(t.commander.DriverName()).Table(t)
}

func (t *Table) GetSelectFields() []string {
	var fields []string
	for _, entity := range t.fields.Get() {
		fields = append(fields, entity.GetValue())
	}

	return fields
}

func (t *Table) ToSqlTable(name string) *Table {
	return NewSqlTable(name, t.Sql())
}

func (t *Table) prepare() {
	t.commander.SetTable(t.name, t.asSql)
	t.commander.SetSelect(t.fields)
	t.commander.SetFrom(t.join)
	t.commander.SetWhere(t.where)
	t.commander.SetGroupBy(t.groupBy)
	t.commander.SetOrderBy(t.orderBy)
	t.commander.SetHaving(t.having)
	t.commander.SetLimit(t.limit)
	t.commander.SetSet(t.set)
}
