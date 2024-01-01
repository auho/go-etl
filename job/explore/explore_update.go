package explore

import (
	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.UpdateModer = (*Update)(nil)

type Update struct {
	*Explore
}

func newUpdateFromExplore(e *Explore) *Update {
	return NewUpdate(e.collect, e.search, e.expression)
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
