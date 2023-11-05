package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/job/mode"
)

func CleanTask(source job.Source, target job.Target, modes []mode.UpdateModer) {
	cleanAction := action.NewClean(target, modes)
	RunTask(source, []action.Actor{cleanAction})
}
