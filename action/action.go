package action

import (
	"github.com/auho/go-toolkit/flow/storage"
	"github.com/auho/go-toolkit/flow/task"
)

const batchSize = 2000

type Actionor interface {
	task.Tasker[storage.MapEntry]
	GetFields() []string
}

type action struct {
	task.Task
}
