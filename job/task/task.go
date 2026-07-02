package task

import (
	"runtime"

	"github.com/auho/go-toolkit-flow/v3/processor/consumer"
	"github.com/auho/go-toolkit-flow/v3/processor/producer"
	"github.com/auho/go-toolkit-flow/v3/storage"
)

const batchSize = 2000

type processor interface {
	GetFields() []string
}

type itemConsumer interface {
	processor

	consumer.Item[storage.MapEntry]
}

type itemProducer interface {
	processor
	producer.Item[storage.MapEntry, storage.MapEntry]
	destinations() ([]storage.Destination[storage.MapEntry], error)
}

type task struct{}

func (t *task) Concurrency() int {
	return runtime.NumCPU()
}

type consumerTask struct {
	task

	consumer.Processor
}

type producerTask struct {
	task

	producer.Processor
}
