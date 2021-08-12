package command

import (
	"fmt"
	"strings"
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
