package explore

import (
	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.InsertModer = (*Insert)(nil)

type Insert struct {
	*Explore
}

func newInsertFromExplore(e *Explore) *Insert {
	return NewInsert(e.collect, e.search, e.expression)
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
	if !token.IsOk() {
		return nil
	}

	ret := token.ToToken()

	i.AddAmount(int64(len(ret)))

	return ret
}
