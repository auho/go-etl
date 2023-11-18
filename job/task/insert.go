package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/job/mode"
)

func InsertTask(source job.Source, target job.Target, moder mode.InsertModer, opts ...func(*action.Insert)) {
	a := action.NewInsert(target, moder, opts...)
	RunTask(source, []action.Actor{a})
}
