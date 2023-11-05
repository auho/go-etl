package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/job/mode"
)

func InsertFlow(source job.Source, target job.Target, moder mode.InsertModer, extraKeys []string) {
	a := action.NewInsert(target, moder, extraKeys)
	RunTask(source, []action.Actor{a})
}
