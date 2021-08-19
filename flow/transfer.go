package flow

import (
	goetl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-simple-db/simple"
)

func TransferFlow(db simple.Driver, config goetl.DbConfig, dataName string, idName string, targetTableName string, alias map[string]string, fixedData map[string]interface{}) {
	transferAction := action.NewTransfer(db, config, targetTableName, alias, fixedData)
	RunFlow(config, dataName, idName, []action.Actionor{transferAction})
}
