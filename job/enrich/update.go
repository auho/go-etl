package enrich

import (
	"github.com/auho/go-etl/v2/job/enrich/collect"
	"github.com/auho/go-etl/v2/job/enrich/condition"
	"github.com/auho/go-etl/v2/job/enrich/search"
	"github.com/auho/go-etl/v2/job/transform"
)

var _ mode.UpdateOperator = (*Update)(nil)

type Update struct {
	*Explore
}

func newUpdateFromExplore(e *Explore) *Update {
	return NewUpdate(e.collect, e.search, e.condition)
}

func NewUpdate(collect collect.Collector, search search.Searcher, expression condition.Operation) *Update {
	return &Update{
		Explore: newExplore(collect, search, expression),
	}
}

func (u *Update) Do(item map[string]any) map[string]any {
	u.AddTotal(1)

	if !u.expressionOperation(item) {
		return nil
	}

	token := u.collect.Do(item, u.search)
	if !token.IsOk() {
		return nil
	}

	ret := token.ToToken()

	u.AddAmount(int64(len(ret)))

	return ret[0]
}
