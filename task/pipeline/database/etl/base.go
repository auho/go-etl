package etl

import (
	"runtime"

	"github.com/auho/go-toolkit-flow/v3/storage/database/destination"
)

// batchSize is the default write batch size.
const batchSize = 2000

// TaskBase provides default Concurrency.
type TaskBase struct{}

func (t *TaskBase) Concurrency() int { return runtime.NumCPU() }

// BaseConfig is the shared config embedded by processor-specific configs.
type BaseConfig struct {
	BatchSize   int
	Concurrency int
}

func (bc *BaseConfig) Check() {
	if bc.BatchSize <= 0 {
		bc.BatchSize = batchSize
	}
	if bc.Concurrency <= 0 {
		bc.Concurrency = runtime.NumCPU()
	}
}

// BulkConfig derives a go-toolkit-flow destination.BulkConfig from BaseConfig.
func (bc *BaseConfig) BulkConfig(isTruncate bool) destination.BulkConfig {
	return destination.BulkConfig{
		IsTruncate:  isTruncate,
		Concurrency: bc.Concurrency,
		BatchSize:   int64(bc.BatchSize),
	}
}
