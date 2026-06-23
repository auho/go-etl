package dml

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command/mysql"
)

const DriverMySQL = "mysql"

func newTableJoinCommand(_driver string) command.TableJoinCommander {
	if _driver == "" {
		_driver = DriverMySQL
	}

	switch _driver {
	case DriverMySQL:
		return mysql.NewTableJoinCommand()
	default:
		panic(fmt.Sprintf("_driver[%s] is not exists", _driver))
	}
}

func newTableCommand(_driver string) command.TableCommander {
	if _driver == "" {
		_driver = DriverMySQL
	}

	switch _driver {
	case DriverMySQL:
		return mysql.NewTableCommand()
	default:
		panic(fmt.Sprintf("driver[%s] is not exists", _driver))
	}
}
