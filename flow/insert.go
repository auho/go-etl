package flow

import (
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

func InsertFlow(db *go_simple_db.SimpleDB, dataName string, idName string, tagTableName string, moder mode.InsertModer, affixFields []string) {
	a := action.NewInsert(db, tagTableName, moder, affixFields)
	RunFlow(db, dataName, idName, []action.Actionor{a})
}
