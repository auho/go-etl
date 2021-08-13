package mysql

import (
	"fmt"
	"strings"
)

type mysqlCommand struct {
}

func NewMysqlCommand() *mysqlCommand {
	return &mysqlCommand{}
}

func (m *mysqlCommand) SelectToString(s []string) string {
	return fmt.Sprintf("SELECT %s ", strings.Join(s, ", "))
}

func (m *mysqlCommand) FromToString(f []string) string {
	return fmt.Sprintf("%s ", strings.Join(f, " "))
}

func (m *mysqlCommand) WhereToString(w []string) string {
	if w == nil {
		return ""
	}

	return fmt.Sprintf("WHERE %s ", strings.Join(w, " AND "))
}

func (m *mysqlCommand) GroupByToString(g []string) string {
	if g == nil {
		return ""
	}

	return fmt.Sprintf("GROUP BY %s ", strings.Join(g, ", "))
}

func (m *mysqlCommand) OrderByToString(o []string) string {
	if o == nil {
		return ""
	}

	return fmt.Sprintf("ORDER BY %s ", strings.Join(o, ", "))
}

func (m *mysqlCommand) LimitToString(l []int) string {
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

func (m *mysqlCommand) addBackQuote(s string) string {
	return fmt.Sprintf("`%s`", s)
}
