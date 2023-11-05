package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/job/mode"
)

func TransferFlow(source job.Source, target job.Target, moder mode.TransferModer) {
	transferAction := action.NewTransfer(target, moder)
	RunTask(source, []action.Actor{transferAction})
}
