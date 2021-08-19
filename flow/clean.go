package flow

import (
	goetl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
	"github.com/auho/go-simple-db/simple"
)

func CleanFlow(db simple.Driver, config goetl.DbConfig, dataName string, idName string, targetTableName string, modes []mode.UpdateModer) {
	cleanAction := action.NewClean(db, config, targetTableName, modes)
	RunFlow(config, dataName, idName, []action.Actionor{cleanAction})
}
