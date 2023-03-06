package flow

import (
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

func CleanFlow(db *goSimpleDb.SimpleDB, dataTable string, idName string, targetTable string, modes []mode.UpdateModer) {
	cleanAction := action.NewClean(db, targetTable, modes)
	RunFlow(db, dataTable, idName, []action.Actioner{cleanAction})
}
