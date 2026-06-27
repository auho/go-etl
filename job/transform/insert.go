package transform

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

var _ InsertOperator = (*Insert)(nil)

type Insert struct {
	*Pipeline
}

func newInsertFromPipeline(p *Pipeline) *Insert {
	return NewInsert(p.collector, p.extractor, p.predicate)
}

func NewInsert(c collect.Collector, e extract.Extractor, p filter.Predicate) *Insert {
	return &Insert{
		Pipeline: newPipeline(c, e, p),
	}
}

func (i *Insert) Apply(item map[string]any) ([]map[string]any, error) {
	i.AddTotal(1)

	ok, err := i.expressionOperation(item)
	if err != nil {
		return nil, fmt.Errorf("expressionOperation: %w", err)
	}
	if !ok {
		return nil, nil
	}

	token, err := i.collector.Search(item, i.extractor)
	if err != nil {
		return nil, fmt.Errorf("collect.Search: %w", err)
	}

	if !token.IsOK() {
		return nil, nil
	}

	ret := token.Rows()

	i.AddAmount(int64(len(ret)))

	return ret, nil
}
