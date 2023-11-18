package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/job/mode"
)

func UpdateAndTransferTask(source job.Source, target job.Target, modes []mode.UpdateModer) {
	a := action.NewUpdateAndTransfer(source, target, modes)
	RunTask(source, []action.Actor{a})
}

func UpdateTask(source job.Source, modes []mode.UpdateModer) {
	a := action.NewUpdate(source, modes)
	RunTask(source, []action.Actor{a})
}
