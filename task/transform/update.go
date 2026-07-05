package transform

import (
	"github.com/auho/go-etl/v3/task/transform/collector"
	"github.com/auho/go-etl/v3/task/transform/filter"
)

var _ UpdateOperator = (*Update)(nil)

type Update struct {
	pipeline
}

func newUpdateFromPipeline(p *pipeline) *Update {
	return NewUpdate(p.collector, p.predicate)
}

func NewUpdate(c *collector.Collector, p filter.Predicate) *Update {
	return &Update{
		pipeline: newPipeline(c, p),
	}
}

func (u *Update) Apply(item map[string]any) (map[string]any, error) {
	ret, err := u.apply(item)
	if err != nil || len(ret) == 0 {
		return nil, err
	}

	return ret[0], nil
}