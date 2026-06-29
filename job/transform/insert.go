package transform

import (
	"github.com/auho/go-etl/v3/job/transform/collector"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

var _ InsertOperator = (*Insert)(nil)

type Insert struct {
	pipeline
}

func newInsertFromPipeline(p *pipeline) *Insert {
	return NewInsert(p.collector, p.predicate)
}

func NewInsert(c *collector.Collector, p filter.Predicate) *Insert {
	return &Insert{
		pipeline: newPipeline(c, p),
	}
}

func (i *Insert) Apply(item map[string]any) ([]map[string]any, error) {
	return i.apply(item)
}