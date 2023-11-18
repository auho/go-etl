package mysql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
	"github.com/auho/go-simple-db/v2/driver/driver"
)

var _ command.TableCommander = (*TableCommand)(nil)
var _ command.Query = (*TableCommand)(nil)

type TableCommand struct {
	mysql
	name          string
	nameBackQuote string
	fields        *command.Entities
	where         string
	groupBy       *command.Entities
	orderBy       *command.Entities
	limit         []int
	join          *command.Join
	set           []*command.Set
	asSql         string
}

func NewTableCommand() *TableCommand {
	return &TableCommand{}
}

func (c *TableCommand) DriverName() string {
	return driver.Mysql
}

func (c *TableCommand) Name() string {
	return c.name
}

func (c *TableCommand) BuildFieldsForInsert() []string {
	s := make([]string, 0)
	for _, field := range c.fields.Get() {
		s = append(s, field.GetValue())
	}

	return s
}

func (c *TableCommand) SetTable(name string, sql string) {
	c.name = name
	c.asSql = sql
	c.nameBackQuote = c.addBackQuote(c.name)
}

func (c *TableCommand) SetSelect(f *command.Entities) {
	c.fields = f
}

func (c *TableCommand) Select() string {
	fs := c.BuildSelect()

	return c.SelectToString(fs)
}

func (c *TableCommand) BuildSelect() []string {
	fields := make([]string, 0)

	for _, v := range c.fields.Get() {
		key, value := v.Get()
		newKey := ""
		if v.IsExpression() {
			newKey = c.addSelfTablePrefix(key)
			fields = append(fields, fmt.Sprintf("%s AS '%s'", newKey, value))
		} else {
			if key == "*" {
				newKey = key
			} else {
				newKey = c.addSelfTablePrefix(c.addBackQuote(key))
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

func (c *TableCommand) SetFrom(j *command.Join) {
	c.join = j
}

func (c *TableCommand) From() string {
	fs := c.BuildFrom()

	return c.FromToString(fs)
}

func (c *TableCommand) BuildFrom() []string {
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
		for k := range c.join.LFields {
			ons = append(ons, fmt.Sprintf("%s.%s = %s.%s",
				c.addBackQuote(c.join.LTable),
				c.addBackQuote(c.join.LFields[k]),
				c.addBackQuote(c.join.RTable),
				c.addBackQuote(c.join.RFields[k]),
			))
		}

		f = f + strings.Join(ons, " AND ")
	}

	return []string{f}
}

func (c *TableCommand) SetWhere(s string) {
	c.where = s
}

func (c *TableCommand) Where() string {
	w := c.BuildWhere()

	return c.WhereToString(w)
}

func (c *TableCommand) BuildWhere() []string {
	if c.where == "" {
		return nil
	}

	return []string{c.addSelfTablePrefix(c.where)}
}

func (c *TableCommand) SetGroupBy(g *command.Entities) {
	c.groupBy = g
}

func (c *TableCommand) GroupBy() string {
	gs := c.BuildGroupBy()

	return c.GroupByToString(gs)
}

func (c *TableCommand) BuildGroupBy() []string {
	if c.groupBy.Len() == 0 {
		return nil
	}

	gs := make([]string, 0)

	for _, v := range c.groupBy.Get() {
		gs = append(gs, c.addSelfTablePrefix(c.addBackQuote(v.GetValue())))
	}

	return gs
}

func (c *TableCommand) SetOrderBy(o *command.Entities) {
	c.orderBy = o
}

func (c *TableCommand) OrderBy() string {
	os := c.BuildOrderBy()

	return c.OrderByToString(os)
}

func (c *TableCommand) BuildOrderBy() []string {
	if c.orderBy.Len() == 0 {
		return nil
	}

	os := make([]string, 0)
	for _, v := range c.orderBy.Get() {
		os = append(os, fmt.Sprintf("%s %s", c.addSelfTablePrefix(c.addBackQuote(v.GetKey())), v.GetValue()))
	}

	return os
}

func (c *TableCommand) SetLimit(l []int) {
	c.limit = l
}

func (c *TableCommand) Limit() string {
	return c.LimitToString(c.limit)
}

func (c *TableCommand) SetSet(s []*command.Set) {
	c.set = s
}

func (c *TableCommand) Set() string {
	ss := c.BuildSet()

	return c.SetToString(ss)
}

func (c *TableCommand) BuildSet() []string {
	ss := make([]string, 0)

	for _, set := range c.set {
		lValue := ""
		rValue := ""
		for k, v := range set.LFields {
			lValue = c.addTablePrefix(c.addBackQuote(v), set.LTable)
			rValueAny := set.RValues[k]
			if set.IsExpression() {
				if _rValue, ok := rValueAny.(string); ok {
					rValue = c.addTablePrefix(_rValue, set.RTable)
				} else {
					panic(fmt.Sprintf("field[%s] set expression is not string, expression[%v]", v, rValueAny))
				}
			} else if set.IsValue() {
				switch _rValue := rValueAny.(type) {
				case string:
					rValue = fmt.Sprintf("'%s'", _rValue)
				default:
					rValue = fmt.Sprintf("%v", _rValue)
				}
			} else {
				if _rValue, ok := rValueAny.(string); ok {
					rValue = c.addTablePrefix(c.addBackQuote(_rValue), set.RTable)
				} else {
					panic(fmt.Sprintf("field[%s] set is not string, expression[%v]", v, rValueAny))
				}
			}

			ss = append(ss, fmt.Sprintf("%s = %s", lValue, rValue))
		}
	}

	return ss
}

func (c *TableCommand) Query() string {
	return fmt.Sprintf("%s%s%s%s%s%s",
		c.Select(),
		c.From(),
		c.Where(),
		c.GroupBy(),
		c.OrderBy(),
		c.Limit(),
	)
}

func (c *TableCommand) InsertQuery(name string) string {
	return c.mysql.insert(name, c)
}

func (c *TableCommand) InsertWithFieldsQuery(name string, fields []string) string {
	return c.mysql.insertWithFields(name, fields, c)
}

func (c *TableCommand) UpdateQuery() string {
	return fmt.Sprintf("UPDATE %s %s%s%s%s", c.nameBackQuote, c.Set(), c.Where(), c.OrderBy(), c.Limit())
}

func (c *TableCommand) DeleteQuery() string {
	return fmt.Sprintf("DELETE %s%s%s%s", c.From(), c.Where(), c.OrderBy(), c.Limit())
}

func (c *TableCommand) addSelfTablePrefix(s string) string {
	return c.addTablePrefix(s, c.name)
}

func (c *TableCommand) addTablePrefix(s string, name string) string {
	name = c.addBackQuote(name)

	re := regexp.MustCompile("([^.])(`[^`.]+`)([^.])")
	s = re.ReplaceAllString(" "+s+" ", fmt.Sprintf("$1%s.$2$3", name))
	return strings.Trim(s, " ")
}
