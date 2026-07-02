package task

import (
	"runtime"
)

type ConfigOption func(*Config)

type Config struct {
	source SourceConfig
}

func (c *Config) Init() {
	c.source.check()
}

type SourceConfig struct {
	Concurrency int
	Maximum     int64
	PageSize    int64
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

type baseConfig struct {
	BatchSize   int
	Concurrency int
}

func (bc *baseConfig) check() {
	if bc.BatchSize <= 0 {
		bc.BatchSize = batchSize
	}

	if bc.Concurrency <= 0 {
		bc.Concurrency = runtime.NumCPU()
	}
}

func WithSourceConfig(sc SourceConfig) func(config *Config) {
	return func(config *Config) {
		config.source = sc
	}
}
