package action

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
)

const batchSize = 2000

type Actor interface {
	task.Singleton[storage.MapEntry]
	GetFields() []string
}

type action struct {
	task.Task

	target job.Target
}
