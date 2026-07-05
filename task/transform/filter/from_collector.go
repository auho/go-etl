package filter

import (
	"fmt"

	"github.com/auho/go-etl/v3/task/transform/collector"
)

var _ Predicate = (*FromCollector)(nil)

// FromCollector adapts a Collector into a Predicate. Match delegates to
// Collector.Extract and reports whether the result is OK; Prepare delegates to
// Collector.Prepare so its error propagates up the pipeline.
type FromCollector struct {
	collector *collector.Collector
}

// NewFromCollector builds a Predicate backed by the given Collector.
func NewFromCollector(c *collector.Collector) *FromCollector {
	return &FromCollector{collector: c}
}

func (f *FromCollector) Prepare() error {
	return f.collector.Prepare()
}

func (f *FromCollector) Match(m map[string]any) (bool, error) {
	token, err := f.collector.Extract(m)
	if err != nil {
		return false, fmt.Errorf("collect.Extract: %w", err)
	}

	return token.IsOK(), nil
}
