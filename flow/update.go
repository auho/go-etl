package flow

import (
	goetl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
)

func UpdateFlow(config goetl.DbConfig, dataName string, idName string, modes []mode.UpdateModer) {
	a := action.NewUpdate(config, dataName, idName, modes)
	RunFlow(config, dataName, idName, []action.Actionor{a})
}
