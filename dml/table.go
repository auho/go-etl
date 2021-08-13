package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
)

type Table struct {
	commander tableCommander
	name      string
	fields    *command.Entries
	where     string
	groupBy   *command.Entries
	orderBy   *command.Entries
	limit     []int
	join      *command.Join
}

func NewTable(name string) *Table {
	t := &Table{}
	t.name = name
	t.fields = command.NewEntries()
	t.groupBy = command.NewEntries()
	t.orderBy = command.NewEntries()
	t.limit = make([]int, 0)

	t.commander = newTableCommand()
	t.commander.SetName(t.name)

	return t
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
		t.fields.AddEntry(k, v)
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

func (t *Table) OrderBy(o map[string]string) *Table {
	for k, v := range o {
		t.orderBy.AddEntry(k, v)
	}

	return t
}

func (t *Table) Limit(start int, offset int) *Table {
	t.limit = []int{start, offset}

	return t
}

func (t *Table) Aggregation(a map[string]string) *Table {
	for k, v := range a {
		e := command.NewAggregationEntry(k, v)
		t.fields.Add(e)
	}

	return t
}

func (t *Table) LeftJoin(keys []string, joinTable *Table, joinKeys []string) *Table {
	t.join = command.NewLeftJoin(t.GetName(), keys, joinTable.GetName(), joinKeys)

	return t
}

func (t *Table) Sql() string {
	return fmt.Sprintf("%s%s%s%s%s%s",
		t.commander.Select(t.fields),
		t.commander.From(t.join),
		t.commander.Where(t.where),
		t.commander.GroupBy(t.groupBy),
		t.commander.OrderBy(t.orderBy),
		t.commander.Limit(t.limit),
	)
}
