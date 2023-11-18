package task

import (
	"runtime"
)

type Config struct {
	sourceConfig SourceConfig
	targetConfig TargetConfig
}

func (c *Config) Check() {
	c.sourceConfig.check()
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
		sc.PageSize = 2000
	}
}

type TargetConfig struct {
}

func WithSourceConfig(sc SourceConfig) func(config *Config) {
	return func(config *Config) {
		config.sourceConfig = sc
	}
}

func WithTargetConfig(tc TargetConfig) func(config *Config) {
	return func(config *Config) {
		config.targetConfig = tc
	}
}
