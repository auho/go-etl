package task

import (
	"runtime"
)

type ConfigOption func(*Config)

type Config struct {
	source SourceConfig
	target TargetConfig
}

func (c *Config) Check() {
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

type TargetConfig struct{}

func WithSourceConfig(sc SourceConfig) func(config *Config) {
	return func(config *Config) {
		config.source = sc
	}
}

func WithTargetConfig(tc TargetConfig) func(config *Config) {
	return func(config *Config) {
		config.target = tc
	}
}
