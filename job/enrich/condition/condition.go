package condition

import (
	"github.com/auho/go-etl/v3/job/enrich/collect"
	"github.com/auho/go-etl/v3/job/enrich/search"
)

var _ Filter = (*Condition)(nil)

type Filter interface {
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

	return token.IsOK()
}

func (c *Condition) ToOperation() Operation {
	return func(m map[string]any) bool {
		return c.OK(m)
	}
}
