package flow

import (
	"github.com/auho/go-etl/v2/action"
	"github.com/auho/go-etl/v2/mode"
)

func UpdateAndTransferFlow(source action.Source, target action.Target, modes []mode.UpdateModer) {
	a := action.NewUpdateAndTransfer(source, target, modes)
	RunFlow(source, []action.Actor{a})
}

func UpdateFlow(source action.Source, modes []mode.UpdateModer) {
	a := action.NewUpdate(source, modes)
	RunFlow(source, []action.Actor{a})
}
