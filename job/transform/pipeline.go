package transform

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collect"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

type Pipeline struct {
	base

	collect   collect.Collector
	search    extract.Extractor
	condition filter.Predicate

	hasExpression bool
	defaultValues map[string]any
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func newPipeline(collect collect.Collector, search extract.Extractor, operation filter.Predicate) *Pipeline {
	return &Pipeline{
		collect:   collect,
		search:    search,
		condition: operation,
	}
}

func (e *Pipeline) expressionOperation(item map[string]any) bool {
	if !e.hasExpression {
		return true
	}

	return e.condition(item)
}

func (e *Pipeline) Title() string {
	return e.GenTitle(e.collect.Title(), e.search.Title())
}

func (e *Pipeline) GetFields() []string {
	return e.collect.Keys()
}

func (e *Pipeline) Keys() []string {
	return e.search.NewExport().Keys()
}

func (e *Pipeline) DefaultValues() map[string]any {
	return maps.Clone(e.defaultValues)
}

func (e *Pipeline) Prepare() error {
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

func (e *Pipeline) Close() error { return nil }

func (e *Pipeline) State() []string {
	return []string{fmt.Sprintf("%s: %s", e.Title(), e.GenCounter())}
}

func (e *Pipeline) SetCollect(collect collect.Collector) *Pipeline {
	e.collect = collect

	return e
}

func (e *Pipeline) SetSearch(search extract.Extractor) *Pipeline {
	e.search = search

	return e
}

func (e *Pipeline) SetCondition(operation filter.Predicate) *Pipeline {
	e.condition = operation

	return e
}

func (e *Pipeline) ToInsert() *Insert {
	return newInsertFromPipeline(e)
}

func (e *Pipeline) ToUpdate() *Update {
	return newUpdateFromPipeline(e)
}
