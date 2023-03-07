package flow

import (
	"github.com/auho/go-etl/action"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

func TransferFlow(db *goSimpleDb.SimpleDB, dataTable, idName, targetTable string, alias map[string]string, fixedData map[string]interface{}) {
	transferAction := action.NewTransfer(db, targetTable, alias, fixedData)
	RunFlow(db, dataTable, idName, []action.Actor{transferAction})
}
