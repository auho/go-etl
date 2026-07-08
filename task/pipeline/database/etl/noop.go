package etl

import (
	"github.com/auho/go-toolkit-flow/v3/processor/consumer"
	"github.com/auho/go-toolkit-flow/v3/storage"
)

var _ consumer.Item[storage.MapEntry] = (*Noop)(nil)

// Noop is a no-operation consumer for testing the ETL scenario without external
// infrastructure.
type Noop struct {
	consumer.Processor
	TaskBase
}

// NewNoop creates a Noop consumer.
func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) Fields() ([]string, error)                { return nil, nil }
func (n *Noop) Summary() string                          { return "Noop" }
func (n *Noop) Prepare() error                           { return nil }
func (n *Noop) BeforeRun() error                         { return nil }
func (n *Noop) Exec(item storage.MapEntry) (bool, error) { return true, nil }
func (n *Noop) AfterRun() error                          { return nil }
func (n *Noop) AppendState()                             {}
func (n *Noop) Close() error                             { return nil }
