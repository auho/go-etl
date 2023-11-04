package flow

import (
	"github.com/auho/go-etl/v2/action"
	"github.com/auho/go-etl/v2/mode"
)

func CleanFlow(source action.Source, target action.Target, modes []mode.UpdateModer) {
	cleanAction := action.NewClean(target, modes)
	RunFlow(source, []action.Actor{cleanAction})
}
