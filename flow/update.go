package flow

import (
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

func UpdateAndTransferFlow(db *goSimpleDb.SimpleDB, dataTable, idName, transferTable string, modes []mode.UpdateModer) {
	a := action.NewUpdateAndTransfer(db, dataTable, transferTable, idName, modes)
	RunFlow(db, dataTable, idName, []action.Actionor{a})
}

func UpdateFlow(db *goSimpleDb.SimpleDB, dataTable, idName string, modes []mode.UpdateModer) {
	a := action.NewUpdate(db, dataTable, idName, modes)
	RunFlow(db, dataTable, idName, []action.Actionor{a})
}
