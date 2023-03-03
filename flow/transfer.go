package flow

import (
	"github.com/auho/go-etl/action"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

func TransferFlow(db *go_simple_db.SimpleDB, dataName string, idName string, targetTableName string, alias map[string]string, fixedData map[string]interface{}) {
	transferAction := action.NewTransfer(db, targetTableName, alias, fixedData)
	RunFlow(db, dataName, idName, []action.Actionor{transferAction})
}
