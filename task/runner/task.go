package runner

import (
	"runtime"

	"github.com/auho/go-toolkit-flow/v3/processor/consumer"
	"github.com/auho/go-toolkit-flow/v3/processor/producer"
	"github.com/auho/go-toolkit-flow/v3/storage"
)

// batchSize is the default write batch size used when a config does not set one.
const batchSize = 2000

// processor declares the source fields a task needs to read.
type processor interface {
	// Fields returns the column names this processor needs from the source.
	Fields() ([]string, error)
}

// itemConsumer is a processor that consumes items without producing output.
type itemConsumer interface {
	processor

	consumer.Item[storage.MapEntry]
}

// itemProducer is a processor that produces items and owns its destinations.
type itemProducer interface {
	processor

	producer.Item[storage.MapEntry, storage.MapEntry]
	// destinations returns the write destinations owned by this producer.
	destinations() ([]storage.Destination[storage.MapEntry], error)
}

// taskBase is the shared base embedded by all task types, providing default concurrency.
type taskBase struct{}

// Concurrency returns the default concurrency (NumCPU).
func (t *taskBase) Concurrency() int {
	return runtime.NumCPU()
}

// consumerTask is the base for consumer tasks.
type consumerTask struct {
	taskBase

	consumer.Processor
}

// producerTask is the base for producer tasks.
type producerTask struct {
	taskBase

	producer.Processor
}
