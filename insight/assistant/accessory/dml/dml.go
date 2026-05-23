package dml

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command/mysql"
)

const DriverMysql = "mysql"

func newTableJoinCommand(_driver string) command.TableJoinCommander {
	if _driver == "" {
		_driver = DriverMysql
	}

	switch _driver {
	case DriverMysql:
		return mysql.NewTableJoinCommand()
	default:
		panic(fmt.Sprintf("_driver[%s] is not exists", _driver))
	}
}

func newTableCommand(_driver string) command.TableCommander {
	if _driver == "" {
		_driver = DriverMysql
	}

	switch _driver {
	case DriverMysql:
		return mysql.NewTableCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", _driver))
	}
}
