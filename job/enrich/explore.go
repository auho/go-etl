package enrich

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
	"github.com/auho/go-etl/v3/job/extract"
)

type Explore struct {
	base

	collect   collect.Collector
	search    extract.Extractor
	condition filter.Predicate

	hasExpression bool
	defaultValues map[string]any
}

func NewExplore() *Explore {
	return &Explore{}
}

func newExplore(collect collect.Collector, search extract.Extractor, operation filter.Predicate) *Explore {
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

func (e *Explore) Title() string {
	return e.genTitle(e.collect.Title(), e.search.Title())
}

func (e *Explore) GetFields() []string {
	return e.collect.Keys()
}

func (e *Explore) Keys() []string {
	return e.search.NewExport().Keys()
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

	e.defaultValues = e.search.NewExport().DefaultValues()

	return nil
}

func (e *Explore) Close() error { return nil }

func (e *Explore) State() []string {
	return []string{fmt.Sprintf("%s: %s", e.Title(), e.genCounter())}
}

func (e *Explore) SetCollect(collect collect.Collector) *Explore {
	e.collect = collect

	return e
}

func (e *Explore) SetSearch(search extract.Extractor) *Explore {
	e.search = search

	return e
}

func (e *Explore) SetCondition(operation filter.Predicate) *Explore {
	e.condition = operation

	return e
}

func (e *Explore) ToInsert() *Insert {
	return newInsertFromExplore(e)
}

func (e *Explore) ToUpdate() *Update {
	return newUpdateFromExplore(e)
}
