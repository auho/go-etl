package condition

import (
	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ Conditioner = (*Condition)(nil)
var _ Conditioner = (AND)(nil)
var _ Conditioner = (OR)(nil)

type Conditioner interface {
	OK(map[string]any) bool
}

type Condition struct {
	collect collect.Collector
	search  search.Searcher
}

func NewCondition(collect collect.Collector, search search.Searcher) *Condition {
	return &Condition{collect: collect, search: search}
}

func (c *Condition) OK(item map[string]any) bool {
	token := c.collect.Do(item, c.search)

	return token.IsOk()
}

type logicalOperation []Conditioner

type AND logicalOperation

func NewAND(cones ...Conditioner) AND {
	a := AND{}
	a = append(a, cones...)

	return a
}

func (a AND) OK(item map[string]any) bool {
	for _, cond := range a {
		if !cond.OK(item) {
			return false
		}
	}

	return true
}

type OR logicalOperation

func NewOR(cones ...Conditioner) OR {
	o := OR{}
	o = append(o, cones...)

	return o
}

func (o OR) OK(item map[string]any) bool {
	for _, cond := range o {
		if cond.OK(item) {
			return true
		}
	}

	return false
}
