package flow

import (
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

func CleanFlow(db *go_simple_db.SimpleDB, dataName string, idName string, targetTableName string, modes []mode.UpdateModer) {
	cleanAction := action.NewClean(db, targetTableName, modes)
	RunFlow(db, dataName, idName, []action.Actionor{cleanAction})
}
