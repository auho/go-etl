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

func newInsertFromPipeline(e *Pipeline) *Insert {
	return NewInsert(e.collect, e.search, e.condition)
}

func NewInsert(collect collect.Collector, search extract.Extractor, expression filter.Predicate) *Insert {
	return &Insert{
		Pipeline: newPipeline(collect, search, expression),
	}
}

func (i *Insert) Apply(item map[string]any) ([]map[string]any, error) {
	i.AddTotal(1)

	if !i.expressionOperation(item) {
		return nil, nil
	}

	token, err := i.collect.Search(item, i.search)
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
