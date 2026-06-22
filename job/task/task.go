package task

import (
	"runtime"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-toolkit-flow/processor/consumer"
	"github.com/auho/go-toolkit-flow/processor/producer"
	"github.com/auho/go-toolkit-flow/storage"
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
	target job.Target

	producer.Processor
}
