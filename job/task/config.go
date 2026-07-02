package task

import (
	"runtime"
)

// baseConfig is the shared configuration embedded by task-specific configs
// (UpdateConfig, InsertConfig, UpdateTransferConfig, CleanConfig). It controls
// the destination write batch size and concurrency.
type baseConfig struct {
	BatchSize   int // destination write batch size; defaults to batchSize (2000)
	Concurrency int // destination write concurrency; defaults to NumCPU
}

func (bc *baseConfig) check() {
	if bc.BatchSize <= 0 {
		bc.BatchSize = batchSize
	}

	if bc.Concurrency <= 0 {
		bc.Concurrency = runtime.NumCPU()
	}
}

// SourceConfig configures the source reader (pagination + concurrency).
type SourceConfig struct {
	Concurrency int   // source scan concurrency; defaults to NumCPU
	Maximum     int64 // max rows to scan; 0 = unlimited
	PageSize    int64 // rows per page; defaults to batchSize (2000)
}

func (sc *SourceConfig) check() {
	if sc.Concurrency <= 0 {
		sc.Concurrency = runtime.NumCPU()
	}

	if sc.Maximum <= 0 {
		sc.Maximum = 0
	}

	if sc.PageSize <= 0 {
		sc.PageSize = batchSize
	}
}

// RunnerOption configures the Runner.
type RunnerOption func(*Runner)

// WithRunnerSourceConfig returns a RunnerOption that sets the source configuration.
func WithRunnerSourceConfig(sc SourceConfig) RunnerOption {
	return func(r *Runner) {
		r.sourceConfig = sc
	}
}
