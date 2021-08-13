package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/auho/go-etl/dml/command"
)

type TableCommand struct {
	mysql
	name          string
	nameBackQuote string
	asSql         string
}

func NewTableCommand() *TableCommand {
	c := &TableCommand{}

	return c
}

func (c *TableCommand) SetTable(name string, sql string) {
	c.name = name
	c.asSql = sql
	c.nameBackQuote = c.addBackQuote(c.name)
}

func (c *TableCommand) Select(f *command.Entries) string {
	fs := c.BuildSelect(f)

	return c.SelectToString(fs)
}

func (c *TableCommand) BuildSelect(f *command.Entries) []string {
	fields := make([]string, 0)

	if f.Len() == 0 {
		fields = append(fields, fmt.Sprintf("%s.*", c.nameBackQuote))
	} else {
		for _, v := range f.Get() {
			key, value := v.Get()
			newKey := ""
			if v.IsAggregation() {
				newKey = c.addTablePrefix(key)
				fields = append(fields, fmt.Sprintf("%s AS '%s'", newKey, value))
			} else {
				newKey = c.addTablePrefix(c.addBackQuote(key))
				if key == value {
					fields = append(fields, newKey)
				} else {
					fields = append(fields, fmt.Sprintf("%s AS '%s'", newKey, value))
				}
			}

		}
	}

	return fields
}

func (c *TableCommand) From(j *command.Join) string {
	fs := c.BuildFrom(j)

	return c.FromToString(fs)
}

func (c *TableCommand) BuildFrom(j *command.Join) []string {
	f := ""
	from := ""
	if c.asSql == "" {
		from = c.nameBackQuote
	} else {
		from = fmt.Sprintf("(%s) AS %s", c.asSql, c.nameBackQuote)
	}

	if j == nil || j.IsFrom() {
		f = fmt.Sprintf("FROM %s", from)
	} else if j.IsLeft() {
		f = fmt.Sprintf("LEFT JOIN %s ON ", from)

		ons := make([]string, 0)
		for k := range j.LKeys {
			ons = append(ons, fmt.Sprintf("%s.%s = %s.%s",
				c.addBackQuote(j.LTable),
				c.addBackQuote(j.LKeys[k]),
				c.addBackQuote(j.RTable),
				c.addBackQuote(j.RKeys[k]),
			))
		}

		f = f + strings.Join(ons, " AND ")
	}

	return []string{f}
}

func (c *TableCommand) Where(s string) string {
	w := c.BuildWhere(s)

	return c.WhereToString(w)
}

func (c *TableCommand) BuildWhere(s string) []string {
	if s == "" {
		return nil
	}

	return []string{c.addTablePrefix(s)}
}

func (c *TableCommand) GroupBy(g *command.Entries) string {
	gs := c.BuildGroupBy(g)

	return c.GroupByToString(gs)
}

func (c *TableCommand) BuildGroupBy(g *command.Entries) []string {
	if g.Len() == 0 {
		return nil
	}

	gs := make([]string, 0)

	for _, v := range g.Get() {
		gs = append(gs, c.addTablePrefix(c.addBackQuote(v.GetValue())))
	}

	return gs
}

func (c *TableCommand) OrderBy(o *command.Entries) string {
	os := c.BuildOrderBy(o)

	return c.OrderByToString(os)
}

func (c *TableCommand) BuildOrderBy(o *command.Entries) []string {
	if o.Len() == 0 {
		return nil
	}

	os := make([]string, 0)
	for _, v := range o.Get() {
		os = append(os, fmt.Sprintf("%s %s", c.addTablePrefix(c.addBackQuote(v.GetKey())), v.GetValue()))
	}

	return os
}

func (c *TableCommand) Limit(l []int) string {
	return c.LimitToString(l)
}

func (c *TableCommand) addBackQuote(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func (c *TableCommand) addTablePrefix(s string) string {
	re := regexp.MustCompile("([^.])(`[^`.]+`)([^.])")
	s = re.ReplaceAllString(" "+s+" ", fmt.Sprintf("$1%s.$2$3", c.nameBackQuote))
	return strings.Trim(s, " ")
}
