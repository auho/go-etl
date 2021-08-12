package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/auho/go-etl/dml/command"
)

type Command struct {
	name          string
	backQuoteName string
}

func NewCommand(name string) *Command {
	c := &Command{}
	c.name = name
	c.backQuoteName = fmt.Sprintf("`%s`", name)

	return c
}

func (c *Command) Select(fieldMap map[string]string) string {
	fields := make([]string, 0)
	for k, v := range fieldMap {
		if k == v {
			fields = append(fields, k)
		} else {
			fields = append(fields, fmt.Sprintf("%s AS %s", k, v))
		}
	}

	return "SELECT " + strings.Join(fields, ", ") + " "
}

func (c *Command) From(j *command.Join) string {
	f := ""
	if j.IsFrom() {
		f = fmt.Sprintf("FROM %s ", c.name)
	} else if j.IsLeft() {
		f = fmt.Sprintf("LEFT JOIN %s ON ", c.name)

		ons := make([]string, 0)
		for k := range j.LKeys {
			ons = append(ons, fmt.Sprintf("%s.%s = %s.%s ",
				j.LTable.GetName(),
				j.LKeys[k],
				j.RTable.GetName(),
				j.RKeys[k],
			))
		}

		f = f + strings.Join(ons, " AND ")
	}

	return f
}

func (c *Command) Where(s string) string {
	return s
}

func (c *Command) GroupBy(g map[string]string) string {
	ss := make([]string, 0)

	for k, v := range g {
		if k == v {
			ss = append(ss)
		} else {

		}

	}

	return strings.Join(ss, ", ")
}

func (c *Command) OrderBy(map[string]string) string {
	return ""
}

func (c *Command) Limit([]int) string {
	return ""
}

func (c *Command) TableName() string {
	return ""
}

func (c *Command) backQuoteField(field string) string {
	re := regexp.MustCompile("([^.])(`[^`.]+`)([^.])")
	return re.ReplaceAllString(field, fmt.Sprintf("$1%s$2$3", c.backQuoteName))
}
