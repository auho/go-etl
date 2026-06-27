package filter

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
)

var _ Spec = (*Filter)(nil)

type Spec interface {
	OK(map[string]any) (bool, error)
	ToPredicate() Predicate
}

type Filter struct {
	collector collect.Collector
	extractor extract.Extractor
}

func NewFilter(c collect.Collector, e extract.Extractor) Predicate {
	f := &Filter{collector: c, extractor: e}

	return f.ToPredicate()
}

func (f *Filter) OK(item map[string]any) (bool, error) {
	token, err := f.collector.Extract(item, f.extractor)
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
