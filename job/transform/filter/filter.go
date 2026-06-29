package filter

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/transform/collector"
)

var _ Spec = (*Filter)(nil)

type Filter struct {
	collector *collector.Collector
}

func NewFilterPredicate(c *collector.Collector) Predicate {
	f := &Filter{collector: c}

	return f.ToPredicate()
}

func (f *Filter) OK(m map[string]any) (bool, error) {
	token, err := f.collector.Extract(m)
	if err != nil {
		return false, fmt.Errorf("collect.Extract: %w", err)
	}

	return token.IsOK(), nil
}

func (f *Filter) ToPredicate() Predicate {
	return func(m map[string]any) (bool, error) {
		return f.OK(m)
	}
}