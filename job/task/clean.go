package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/action"
	"github.com/auho/go-etl/v2/job/mode"
)

func CleanTask(resource job.CleanResource, modes []mode.UpdateModer, opts ...func(clean *action.Clean)) {
	cleanAction := action.NewClean(resource, modes, opts...)
	RunTask(resource.SourceTarget(), []action.Actor{cleanAction})
}
