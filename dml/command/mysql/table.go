package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/auho/go-etl/dml/command"
)

type tableCommand struct {
	mysql
	name          string
	nameBackQuote string
	fields        *command.Entries
	where         string
	groupBy       *command.Entries
	orderBy       *command.Entries
	limit         []int
	join          *command.Join
	asSql         string
}

func NewTableCommand() *tableCommand {
	return &tableCommand{}
}

func (c *tableCommand) Name() string {
	return c.name
}

func (c *tableCommand) BuildFieldsForInsert() []string {
	s := make([]string, 0)
	for _, field := range c.fields.Get() {
		s = append(s, field.GetValue())
	}

	return s
}

func (c *tableCommand) SetTable(name string, sql string) {
	c.name = name
	c.asSql = sql
	c.nameBackQuote = c.addBackQuote(c.name)
}

func (c *tableCommand) SetSelect(f *command.Entries) {
	c.fields = f
}

func (c *tableCommand) Select() string {
	fs := c.BuildSelect()

	return c.SelectToString(fs)
}

func (c *tableCommand) BuildSelect() []string {
	fields := make([]string, 0)

	for _, v := range c.fields.Get() {
		key, value := v.Get()
		newKey := ""
		if v.IsAggregation() {
			newKey = c.addTablePrefix(key)
			fields = append(fields, fmt.Sprintf("%s AS '%s'", newKey, value))
		} else {
			if key == "*" {
				newKey = key
			} else {
				newKey = c.addTablePrefix(c.addBackQuote(key))
			}

			if key == value {
				fields = append(fields, newKey)
			} else {
				fields = append(fields, fmt.Sprintf("%s AS '%s'", newKey, value))
			}
		}
	}

	return fields
}

func (c *tableCommand) SetFrom(j *command.Join) {
	c.join = j
}

func (c *tableCommand) From() string {
	fs := c.BuildFrom()

	return c.FromToString(fs)
}

func (c *tableCommand) BuildFrom() []string {
	f := ""
	from := ""
	if c.asSql == "" {
		from = c.nameBackQuote
	} else {
		from = fmt.Sprintf("(%s) AS %s", c.asSql, c.nameBackQuote)
	}

	if c.join == nil || c.join.IsFrom() {
		f = fmt.Sprintf("FROM %s", from)
	} else if c.join.IsLeft() {
		f = fmt.Sprintf("LEFT JOIN %s ON ", from)

		ons := make([]string, 0)
		for k := range c.join.LKeys {
			ons = append(ons, fmt.Sprintf("%s.%s = %s.%s",
				c.addBackQuote(c.join.LTable),
				c.addBackQuote(c.join.LKeys[k]),
				c.addBackQuote(c.join.RTable),
				c.addBackQuote(c.join.RKeys[k]),
			))
		}

		f = f + strings.Join(ons, " AND ")
	}

	return []string{f}
}

func (c *tableCommand) SetWhere(s string) {
	c.where = s
}

func (c *tableCommand) Where() string {
	w := c.BuildWhere()

	return c.WhereToString(w)
}

func (c *tableCommand) BuildWhere() []string {
	if c.where == "" {
		return nil
	}

	return []string{c.addTablePrefix(c.where)}
}

func (c *tableCommand) SetGroupBy(g *command.Entries) {
	c.groupBy = g
}

func (c *tableCommand) GroupBy() string {
	gs := c.BuildGroupBy()

	return c.GroupByToString(gs)
}

func (c *tableCommand) BuildGroupBy() []string {
	if c.groupBy.Len() == 0 {
		return nil
	}

	gs := make([]string, 0)

	for _, v := range c.groupBy.Get() {
		gs = append(gs, c.addTablePrefix(c.addBackQuote(v.GetValue())))
	}

	return gs
}

func (c *tableCommand) SetOrderBy(o *command.Entries) {
	c.orderBy = o
}

func (c *tableCommand) OrderBy() string {
	os := c.BuildOrderBy()

	return c.OrderByToString(os)
}

func (c *tableCommand) BuildOrderBy() []string {
	if c.orderBy.Len() == 0 {
		return nil
	}

	os := make([]string, 0)
	for _, v := range c.orderBy.Get() {
		os = append(os, fmt.Sprintf("%s %s", c.addTablePrefix(c.addBackQuote(v.GetKey())), v.GetValue()))
	}

	return os
}

func (c *tableCommand) SetLimit(l []int) {
	c.limit = l
}

func (c *tableCommand) Limit() string {
	return c.LimitToString(c.limit)
}

func (c *tableCommand) Query() string {
	return fmt.Sprintf("%s%s%s%s%s%s",
		c.Select(),
		c.From(),
		c.Where(),
		c.GroupBy(),
		c.OrderBy(),
		c.Limit(),
	)
}

func (c *tableCommand) InsertQuery(name string) string {
	return c.mysql.insert(name, c)
}

func (c *tableCommand) InsertWithFieldsQuery(name string, fields []string) string {
	return c.mysql.insertWithFields(name, fields, c)
}

func (c *tableCommand) DeleteQuery() string {
	return fmt.Sprintf("DELETE %s%s%s%s", c.From(), c.Where(), c.OrderBy(), c.Limit())
}

func (c *tableCommand) addTablePrefix(s string) string {
	re := regexp.MustCompile("([^.])(`[^`.]+`)([^.])")
	s = re.ReplaceAllString(" "+s+" ", fmt.Sprintf("$1%s.$2$3", c.nameBackQuote))
	return strings.Trim(s, " ")
}
