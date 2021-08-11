package flow

import (
	goEtl "github.com/auho/go-etl"
	"github.com/auho/go-etl/action"
	"github.com/auho/go-etl/mode"
)

func UpdateFlow(config goEtl.DbConfig, dataName string, idName string, modes []mode.UpdateModer) {
	a := action.NewUpdate(config, dataName, idName, modes)
	RunFlow(config, dataName, idName, []action.Actionor{a})
}
