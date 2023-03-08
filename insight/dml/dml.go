package dml

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/dml/command"
	"github.com/auho/go-etl/v2/insight/dml/command/mysql"
)

const DriverMysql = "mysql"

var driver = ""

func RegisterDriver(d string) {
	driver = d
}

func newTableJoinCommand() command.TableJoinCommander {
	switch driver {
	case DriverMysql:
		return mysql.NewTableJoinCommand()
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
