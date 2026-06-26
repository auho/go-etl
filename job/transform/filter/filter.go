package filter

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
)

var _ Spec = (*Matcher)(nil)

type Spec interface {
	OK(map[string]any) bool
	ToPredicate() Predicate
}

type Matcher struct {
	collect collect.Collector
	search  extract.Extractor
}

func NewMatcher(collect collect.Collector, search extract.Extractor) Predicate {
	c := &Matcher{collect: collect, search: search}

	return c.ToPredicate()
}

func (c *Matcher) OK(item map[string]any) bool {
	token := c.collect.Search(item, c.search)

	return token.IsOK()
}

func (c *Matcher) ToPredicate() Predicate {
	return func(m map[string]any) bool {
		return c.OK(m)
	}
}
