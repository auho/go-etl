package database

import (
	"fmt"
	"time"
)

type DbSourceConfig struct {
	MaxConcurrent int
	Size          int
	Page          int
	Driver        string
	Dsn           string
	Table         string
	PKeyName      string
	Fields        []string
}

func NewDbSourceConfig() *DbSourceConfig {
	c := &DbSourceConfig{}

	return c
}

func (sc *DbSourceConfig) check() {
	if sc.MaxConcurrent <= 0 {
		panic(fmt.Sprintf("db source config Max Concurrent is error %d", sc.MaxConcurrent))
	}

	if sc.Size <= 0 {
		panic(fmt.Sprintf("db source config Size is error %d", sc.Size))
	}

	if sc.PKeyName == "" {
		panic(fmt.Sprintf("db source config PKeyName is error %s", sc.PKeyName))
	}

	if sc.Table == "" {
		panic(fmt.Sprintf("db source config Table is error %s", sc.Table))
	}
}

type DbTargetConfig struct {
	MaxConcurrent int
	Size          int
	Driver        string
	Dsn           string
	Table         string
}

func NewDbTargetConfig() *DbTargetConfig {
	c := &DbTargetConfig{}

	return c
}

func (tc *DbTargetConfig) check() {
	if tc.MaxConcurrent <= 0 {
		panic(fmt.Sprintf("db target config Max Concurrent is error %d", tc.MaxConcurrent))
	}

	if tc.Size <= 0 {
		panic(fmt.Sprintf("db target config Size is error %d", tc.Size))
	}

	if tc.Table == "" {
		panic(fmt.Sprintf("db target config Table is error %s", tc.Table))
	}
}

type DbSourceState struct {
	duration      time.Duration
	maxConcurrent int
	page          int
	size          int
	itemAmount    int64
}

func newDbSourceState() *DbSourceState {
	s := &DbSourceState{}

	return s
}

type DbTargetState struct {
	duration      time.Duration
	maxConcurrent int
	size          int
	itemAmount    int64
}

func newDbTargetState() *DbTargetState {
	s := &DbTargetState{}

	return s
}
