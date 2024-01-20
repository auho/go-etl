package explore

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/condition"
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/mode"
)

type Explore struct {
	mode.Mode

	collect   collect.Collector
	search    search.Searcher
	condition condition.Operation

	hasExpression bool
	defaultValues map[string]any
}

func GenExplore() *Explore {
	return &Explore{}
}

func newExplore(collect collect.Collector, search search.Searcher, operation condition.Operation) *Explore {
	return &Explore{
		collect:   collect,
		search:    search,
		condition: operation,
	}
}

func (e *Explore) expressionOperation(item map[string]any) bool {
	if !e.hasExpression {
		return true
	}

	return e.condition(item)
}

func (e *Explore) GetTitle() string {
	return e.GenTitle(e.collect.GetTitle(), e.search.GetTitle())
}

func (e *Explore) GetFields() []string {
	return e.collect.GetKeys()
}

func (e *Explore) GetKeys() []string {
	return nil // TODO
}

func (e *Explore) DefaultValues() map[string]any {
	return maps.Clone(e.defaultValues)
}

func (e *Explore) Prepare() error {
	err := e.search.Prepare()
	if err != nil {
		return err
	}

	if e.condition != nil {
		e.hasExpression = true
	}

	e.defaultValues = e.search.GenExport().GetDefaultValues()

	return nil
}

func (e *Explore) Close() error { return nil }

func (e *Explore) State() []string {
	return []string{fmt.Sprintf("%s: %s", e.GetTitle(), e.GenCounter())}
}

func (e *Explore) SetCollect(collect collect.Collector) *Explore {
	e.collect = collect

	return e
}

func (e *Explore) SetSearch(search search.Searcher) *Explore {
	e.search = search

	return e
}

func (e *Explore) SetCondition(operation condition.Operation) *Explore {
	e.condition = operation

	return e
}

func (e *Explore) ToInsert() *Insert {
	return newInsertFromExplore(e)
}

func (e *Explore) ToUpdate() *Update {
	return newUpdateFromExplore(e)
}
