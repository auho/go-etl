package flow

import (
	"github.com/auho/go-etl/v2/action"
	"github.com/auho/go-etl/v2/mode"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

func InsertFlow(db *goSimpleDb.SimpleDB, dataTable, idName, targetTable string, moder mode.InsertModer, affixFields []string) {
	a := action.NewInsert(db, targetTable, moder, affixFields)
	RunFlow(db, dataTable, idName, []action.Actor{a})
}
