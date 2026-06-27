package transform

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

var _ UpdateOperator = (*Update)(nil)

type Update struct {
	pipeline
}

func newUpdateFromPipeline(p *pipeline) *Update {
	return NewUpdate(p.collector, p.extractor, p.predicate)
}

func NewUpdate(c collect.Collector, e extract.Extractor, p filter.Predicate) *Update {
	return &Update{
		pipeline: newPipeline(c, e, p),
	}
}

func (u *Update) Apply(item map[string]any) (map[string]any, error) {
	ret, err := u.apply(item)
	if err != nil || len(ret) == 0 {
		return nil, err
	}

	return ret[0], nil
}
