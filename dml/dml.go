package dml

import (
	"fmt"

	"github.com/auho/go-etl/dml/command"
	"github.com/auho/go-etl/dml/command/mysql"
)

const DriverMysql = "mysql"

const reservedSelect = "select"
const reservedFrom = "from"
const reservedWhere = "where"
const reservedGroupBy = "groupBy"
const reservedOrderBy = "orderBy"
const reservedLimit = "limit"

var driver = ""

func RegisterDriver(d string) {
	driver = d
}

func newDriverCommand() command.DriverCommander {
	switch driver {
	case DriverMysql:
		return mysql.NewMysqlCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", driver))
	}
}

func newTableCommand() command.TableCommander {
	switch driver {
	case DriverMysql:
		return mysql.NewTableCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", driver))
	}
}

func newInsertCommand() command.InsertCommander {
	switch driver {
	case DriverMysql:
		return mysql.NewInsertCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", driver))
	}
}
