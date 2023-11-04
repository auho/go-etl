package flow

import (
	"github.com/auho/go-etl/v2/action"
)

func TransferFlow(source action.Source, target action.Target, keys []string, alias map[string]string, fixedData map[string]any) {
	transferAction := action.NewTransfer(target, keys, alias, fixedData)
	RunFlow(source, []action.Actor{transferAction})
}
