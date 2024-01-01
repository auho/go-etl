package expression

import (
	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ Conditioner = (*Condition)(nil)
var _ Conditioner = (AND)(nil)
var _ Conditioner = (OR)(nil)

type Conditioner interface {
	OK(map[string]any) bool
	ToOperation() Operation
}

type Condition struct {
	collect collect.Collector
	search  search.Searcher
}

func NewCondition(collect collect.Collector, search search.Searcher) Operation {
	c := &Condition{collect: collect, search: search}

	return c.ToOperation()
}

func (c *Condition) OK(item map[string]any) bool {
	token := c.collect.Do(item, c.search)

	return token.IsOk()
}

func (c *Condition) ToOperation() Operation {
	return func(m map[string]any) bool {
		return c.OK(m)
	}
}
