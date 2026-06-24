package enrich

import (
	"github.com/auho/go-etl/v3/job/enrich/collect"
	"github.com/auho/go-etl/v3/job/enrich/condition"
	"github.com/auho/go-etl/v3/job/enrich/search"
	"github.com/auho/go-etl/v3/job/transform"
)

var _ transform.InsertOperator = (*Insert)(nil)

type Insert struct {
	*Explore
}

func newInsertFromExplore(e *Explore) *Insert {
	return NewInsert(e.collect, e.search, e.condition)
}

func NewInsert(collect collect.Collector, search search.Searcher, expression condition.Operation) *Insert {
	return &Insert{
		Explore: newExplore(collect, search, expression),
	}
}

func (i *Insert) Do(item map[string]any) []map[string]any {
	i.AddTotal(1)

	if !i.expressionOperation(item) {
		return nil
	}

	token := i.collect.Do(item, i.search)
	if !token.IsOK() {
		return nil
	}

	ret := token.ToToken()

	i.AddAmount(int64(len(ret)))

	return ret
}
