package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
	"github.com/auho/go-etl/dml/command/mysql"
)

const DriverMysql = "mysql"

var driver = ""

type tableCommander interface {
	SetName(string)
	Select(*command.SortMap) string
	BuildSelect(*command.SortMap) []string
	From(join *command.Join) string
	BuildFrom(j *command.Join) []string
	Where(string) string
	BuildWhere(s string) []string
	GroupBy(*command.SortMap) string
	BuildGroupBy(map[string]string) []string
	OrderBy(map[string]string) string
	BuildOrderBy(map[string]string) []string
	Limit([]int) string
}

type driverCommander interface {
	SelectToString([]string) string
	FromToString([]string) string
	WhereToString([]string) string
	GroupByToString([]string) string
	OrderByToString([]string) string
	LimitToString([]int) string
}

func RegisterDriver(d string) {
	driver = d
}

func newTableCommand() tableCommander {
	switch driver {
	case DriverMysql:
		return mysql.NewTableCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", driver))
	}
}

func newDriverCommand() driverCommander {
	switch driver {
	case DriverMysql:
		return mysql.NewMysqlCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", driver))
	}
}
