package command

import (
	"fmt"
)

type Table struct {
	commander tableCommander
	name      string
	fields    map[string]string
	where     string
	groupBy   map[string]string
	orderBy   map[string]string
	limit     []int
	join      *Join
}

func NewTable(name string) *Table {
	t := &Table{}
	t.name = name

	return t
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) Select(fields []string) *Table {
	for _, field := range fields {
		t.fields[field] = field
	}

	return t
}

func (t *Table) SelectAlias(alias map[string]string) *Table {
	for k, v := range alias {
		t.fields[k] = v
	}

	return t
}

func (t *Table) Where(s string) *Table {
	t.where = s

	return t
}

func (t *Table) GroupBy(g []string) *Table {
	for _, v := range g {
		t.groupBy[v] = v
		t.fields[v] = v
	}

	return t
}

func (t *Table) GroupByAlias(g map[string]string) *Table {
	for k, v := range g {
		t.groupBy[k] = k
		t.fields[k] = v
	}

	return t
}

func (t *Table) OrderBy(o map[string]string) *Table {
	for k, v := range o {
		t.orderBy[k] = v
	}

	return t
}

func (t *Table) Limit(start int, offset int) *Table {
	t.limit = []int{start, offset}

	return t
}

func (t *Table) Aggregation(a map[string]string) *Table {
	for k, v := range a {
		t.fields[k] = v
	}

	return t
}

func (t *Table) LeftJoin(keys []string, joinTable *Table, joinKeys []string) *Table {
	t.join = newLeftJoin(t, keys, joinTable, joinKeys)

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
