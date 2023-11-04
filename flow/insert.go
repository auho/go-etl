package flow

import (
	"github.com/auho/go-etl/v2/action"
	"github.com/auho/go-etl/v2/mode"
)

func InsertFlow(source action.Source, target action.Target, moder mode.InsertModer, extraKeys []string) {
	a := action.NewInsert(target, moder, extraKeys)
	RunFlow(source, []action.Actor{a})
}
