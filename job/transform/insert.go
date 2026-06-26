package transform

import (
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

func (i *Insert) Apply(item map[string]any) []map[string]any {
	i.AddTotal(1)

	if !i.expressionOperation(item) {
		return nil
	}

	token := i.collect.Search(item, i.search)
	if !token.IsOK() {
		return nil
	}

	ret := token.Rows()

	i.AddAmount(int64(len(ret)))

	return ret
}
