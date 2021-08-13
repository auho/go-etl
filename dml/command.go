package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
	"github.com/auho/go-etl/dml/command/mysql"
)

const DriverMysql = "mysql"

const reservedSelect = "select"
const reservedFrom = "from"
const reservedWhere = "wherer"
const reservedGroupBy = "groupBy"
const reservedOrderBy = "orderBy"
const reservedLimit = "limit"

var driver = ""

type tableCommander interface {
	SetTable(string, string)
	Select(*command.Entries) string
	BuildSelect(*command.Entries) []string
	From(*command.Join) string
	BuildFrom(j *command.Join) []string
	Where(string) string
	BuildWhere(s string) []string
	GroupBy(*command.Entries) string
	BuildGroupBy(*command.Entries) []string
	OrderBy(*command.Entries) string
	BuildOrderBy(*command.Entries) []string
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
