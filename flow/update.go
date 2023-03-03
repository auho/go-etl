package flow

import (
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

func UpdateFlow(config *go_simple_db.SimpleDB, dataName string, idName string, modes []mode.UpdateModer) {
	a := action.NewUpdate(config, dataName, idName, modes)
	RunFlow(config, dataName, idName, []action.Actionor{a})
}
