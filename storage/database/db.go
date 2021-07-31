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

type DbState struct {
	duration      time.Duration
	maxConcurrent int
	size          int
	itemAmount    int64
	status        string
	title         string
}

func (s *DbState) GetTitle() string {
	return s.title
}

func (s *DbState) GetStatus() string {
	return s.status
}

type DbSourceState struct {
	DbState
	page int
}

func newDbSourceState() *DbSourceState {
	s := &DbSourceState{}
	s.title = "source"

	return s
}

func (s *DbSourceState) DoneStatus() string {
	return fmt.Sprintf("Max Concurrent: %d Page: %d Size: %d Amount: %d", s.maxConcurrent, s.page, s.size, s.itemAmount)
}

type DbTargetState struct {
	DbState
}

func newDbTargetState() *DbTargetState {
	s := &DbTargetState{}
	s.title = "target"

	return s
}

func (s *DbTargetState) DoneStatus() string {
	return fmt.Sprintf("Max Concurrent: %d Size: %d Amount: %d", s.maxConcurrent, s.size, s.itemAmount)
}
