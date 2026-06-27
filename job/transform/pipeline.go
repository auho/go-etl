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

	collector collect.Collector
	extractor extract.Extractor
	predicate filter.Predicate

	hasExpression bool
	defaultValues map[string]any
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func newPipeline(c collect.Collector, e extract.Extractor, p filter.Predicate) *Pipeline {
	return &Pipeline{
		collector: c,
		extractor: e,
		predicate: p,
	}
}

func (p *Pipeline) expressionOperation(item map[string]any) (bool, error) {
	if !p.hasExpression {
		return true, nil
	}

	return p.predicate(item)
}

func (p *Pipeline) Title() string {
	return p.GenTitle(p.collector.Title(), p.extractor.Title())
}

func (p *Pipeline) GetFields() []string {
	return p.collector.Keys()
}

func (p *Pipeline) Keys() []string {
	return p.extractor.NewExport().Keys()
}

func (p *Pipeline) DefaultValues() map[string]any {
	return maps.Clone(p.defaultValues)
}

func (p *Pipeline) Prepare() error {
	err := p.extractor.Prepare()
	if err != nil {
		return err
	}

	if p.predicate != nil {
		p.hasExpression = true
	}

	p.defaultValues = p.extractor.NewExport().DefaultValues()

	return nil
}

func (p *Pipeline) Close() error { return nil }

func (p *Pipeline) State() []string {
	return []string{fmt.Sprintf("%s: %s", p.Title(), p.GenCounter())}
}

func (p *Pipeline) SetCollector(c collect.Collector) *Pipeline {
	p.collector = c

	return p
}

func (p *Pipeline) SetExtractor(e extract.Extractor) *Pipeline {
	p.extractor = e

	return p
}

func (p *Pipeline) SetPredicate(predicate filter.Predicate) *Pipeline {
	p.predicate = predicate

	return p
}

func (p *Pipeline) ToInsert() *Insert {
	return newInsertFromPipeline(p)
}

func (p *Pipeline) ToUpdate() *Update {
	return newUpdateFromPipeline(p)
}
