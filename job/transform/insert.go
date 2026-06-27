package transform

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

var _ InsertOperator = (*Insert)(nil)

type Insert struct {
	pipeline
}

func newInsertFromPipeline(p *pipeline) *Insert {
	return NewInsert(p.collector, p.extractor, p.predicate)
}

func NewInsert(c collect.Collector, e extract.Extractor, p filter.Predicate) *Insert {
	return &Insert{
		pipeline: newPipeline(c, e, p),
	}
}

func (i *Insert) Apply(item map[string]any) ([]map[string]any, error) {
	return i.apply(item)
}
