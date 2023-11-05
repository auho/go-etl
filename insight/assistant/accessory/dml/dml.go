package dml

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
	mysql2 "github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command/mysql"
)

const DriverMysql = "mysql"

var driver = ""

func RegisterDriverMysql() {
	RegisterDriver(DriverMysql)
}

func RegisterDriver(d string) {
	driver = d
}

func newTableJoinCommand() command.TableJoinCommander {
	switch driver {
	case DriverMysql:
		return mysql2.NewTableJoinCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", driver))
	}
}

func newTableCommand() command.TableCommander {
	switch driver {
	case DriverMysql:
		return mysql2.NewTableCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", driver))
	}
}
