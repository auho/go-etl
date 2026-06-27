package transform

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

var _ UpdateOperator = (*Update)(nil)

type Update struct {
	*Pipeline
}

func newUpdateFromPipeline(p *Pipeline) *Update {
	return NewUpdate(p.collector, p.extractor, p.predicate)
}

func NewUpdate(c collect.Collector, e extract.Extractor, p filter.Predicate) *Update {
	return &Update{
		Pipeline: newPipeline(c, e, p),
	}
}

func (u *Update) Apply(item map[string]any) (map[string]any, error) {
	u.addTotal(1)

	ok, err := u.evaluatePredicate(item)
	if err != nil {
		return nil, fmt.Errorf("evaluatePredicate: %w", err)
	}
	if !ok {
		return nil, nil
	}

	token, err := u.collector.Extract(item, u.extractor)
	if err != nil {
		return nil, fmt.Errorf("collect.Extract: %w", err)
	}
	if !token.IsOK() {
		return nil, nil
	}

	ret := token.Rows()

	u.addAmount(int64(len(ret)))

	return ret[0], nil
}
