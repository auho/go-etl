package action

import (
	simpleDb "github.com/auho/go-simple-db/v2"
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
}

type Source interface {
	GetIdName() string
	TableName() string
	GetDB() *simpleDb.SimpleDB
}

type Target interface {
	GetIdName() string
	TableName() string
	GetDB() *simpleDb.SimpleDB
}
