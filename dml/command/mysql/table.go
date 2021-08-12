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
}

func NewTableCommand() *TableCommand {
	c := &TableCommand{}

	return c
}

func (c *TableCommand) SetName(n string) {
	c.name = n
	c.nameBackQuote = fmt.Sprintf("`%s`", c.name)
}

func (c *TableCommand) Select(fieldMap map[string]string) string {
	fields := c.BuildSelect(fieldMap)

	return c.SelectToString(fields)
}

func (c *TableCommand) BuildSelect(fieldMap map[string]string) []string {
	fields := make([]string, 0)

	if len(fieldMap) == 0 {
		fields = append(fields, fmt.Sprintf("%s.*", c.nameBackQuote))
	} else {
		for k, v := range fieldMap {
			k = c.backQuote(k)

			if k == v {
				fields = append(fields, k)
			} else {
				fields = append(fields, fmt.Sprintf("%s AS '%s'", k, v))
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
	if j == nil || j.IsFrom() {
		f = fmt.Sprintf("FROM %s ", c.nameBackQuote)
	} else if j.IsLeft() {
		f = fmt.Sprintf("LEFT JOIN %s ON ", c.nameBackQuote)

		ons := make([]string, 0)
		for k := range j.LKeys {
			ons = append(ons, fmt.Sprintf("%s.%s = %s.%s ",
				c.backQuote(j.LTable),
				c.backQuote(j.LKeys[k]),
				c.backQuote(j.RTable),
				c.backQuote(j.RKeys[k]),
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

	return []string{c.backQuote(s)}
}

func (c *TableCommand) GroupBy(g map[string]string) string {
	gs := c.BuildGroupBy(g)

	return c.GroupByToString(gs)
}

func (c *TableCommand) BuildGroupBy(g map[string]string) []string {
	if len(g) == 0 {
		return nil
	}

	gs := make([]string, 0)

	for k := range g {
		gs = append(gs, c.backQuote(k))
	}

	return gs
}

func (c *TableCommand) OrderBy(o map[string]string) string {
	os := c.BuildOrderBy(o)

	return c.OrderByToString(os)
}

func (c *TableCommand) BuildOrderBy(o map[string]string) []string {
	if len(o) == 0 {
		return nil
	}

	os := make([]string, 0)
	for k, v := range o {
		os = append(os, fmt.Sprintf("%s %s", c.backQuote(k), v))
	}

	return os
}

func (c *TableCommand) Limit(l []int) string {
	return c.LimitToString(l)
}

func (c *TableCommand) backQuote(field string) string {
	re := regexp.MustCompile("([^.])(`[^`.]+`)([^.])")
	return re.ReplaceAllString(field, fmt.Sprintf("$1%s$2$3", c.nameBackQuote))
}
