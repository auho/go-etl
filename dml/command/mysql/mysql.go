package mysql

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/dml/command"
)

type mysql struct {
}

func (m *mysql) SelectToString(s []string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(s, ", "))
}

func (m *mysql) FromToString(f []string) string {
	return fmt.Sprintf("%s ", strings.Join(f, " "))
}

func (m *mysql) WhereToString(w []string) string {
	if w == nil {
		return ""
	}

	return fmt.Sprintf("WHERE %s ", strings.Join(w, " AND "))
}

func (m *mysql) GroupByToString(g []string) string {
	if g == nil {
		return ""
	}

	return fmt.Sprintf("GROUP BY %s ", strings.Join(g, ", "))
}

func (m *mysql) OrderByToString(o []string) string {
	if o == nil {
		return ""
	}

	return fmt.Sprintf("ORDER BY %s ", strings.Join(o, ", "))
}

func (m *mysql) LimitToString(l []int) string {
	if len(l) == 0 {
		return ""
	}

	s := ""
	if len(l) == 1 {
		s = fmt.Sprintf("%d", l[0])
	} else {
		s = fmt.Sprintf("%d, %d", l[0], l[1])
	}

	return "LIMIT " + s
}

func (m *mysql) insert(name string, q command.Query) string {
	return fmt.Sprintf("INSERT INTO %s %s", m.addBackQuote(name), q.Query())
}

func (m *mysql) insertWithFields(name string, fields []string, q command.Query) string {
	s := ""
	if fields == nil {
		fields = q.BuildFieldsForInsert()
	}

	for k, field := range fields {
		fields[k] = m.addBackQuote(field)
	}

	s = fmt.Sprintf(" (%s) ", strings.Join(fields, ", "))

	return fmt.Sprintf("INSERT INTO %s%s %s", m.addBackQuote(name), s, q.Query())
}

func (m *mysql) addBackQuote(s string) string {
	return fmt.Sprintf("`%s`", s)
}
