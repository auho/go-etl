package explore

import (
	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/condition"
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.InsertModer = (*Insert)(nil)

type Insert struct {
	*Explore
}

func NewInsert(collect collect.Collector, search search.Searcher, condition condition.Conditioner) *Insert {
	return &Insert{
		Explore: newExplore(collect, search, condition),
	}
}

func (i *Insert) Do(item map[string]any) []map[string]any {
	i.AddTotal(1)

	if !i.doCondition(item) {
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
