package transform

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v3/job/transform/collector"
	"github.com/auho/go-etl/v3/job/transform/filter"
)

type pipeline struct {
	base

	collector  *collector.Collector
	predicate  filter.Predicate

	hasPredicate  bool
	defaultValues map[string]any
}

func newPipeline(c *collector.Collector, p filter.Predicate) pipeline {
	return pipeline{
		collector: c,
		predicate: p,
	}
}

func (p *pipeline) evaluatePredicate(item map[string]any) (bool, error) {
	if !p.hasPredicate {
		return true, nil
	}

	return p.predicate(item)
}

func (p *pipeline) apply(item map[string]any) ([]map[string]any, error) {
	p.addTotal(1)

	ok, err := p.evaluatePredicate(item)
	if err != nil {
		return nil, fmt.Errorf("evaluatePredicate: %w", err)
	}
	if !ok {
		return nil, nil
	}

	token, err := p.collector.Extract(item)
	if err != nil {
		return nil, fmt.Errorf("collect.Extract: %w", err)
	}
	if !token.IsOK() {
		return nil, nil
	}

	ret := token.Rows()
	p.addAmount(int64(len(ret)))

	return ret, nil
}

func (p *pipeline) Title() string {
	return p.collector.Title()
}

func (p *pipeline) GetFields() []string {
	return p.collector.Fields()
}

func (p *pipeline) Keys() []string {
	return p.collector.Keys()
}

func (p *pipeline) DefaultValues() map[string]any {
	return maps.Clone(p.defaultValues)
}

func (p *pipeline) Prepare() error {
	if err := p.collector.Prepare(); err != nil {
		return err
	}

	if p.predicate != nil {
		p.hasPredicate = true
	}

	p.defaultValues = p.collector.DefaultValues()

	return nil
}

func (p *pipeline) Close() error {
	return p.collector.Close()
}

func (p *pipeline) State() []string {
	return []string{fmt.Sprintf("%s: %s", p.Title(), p.genCounter())}
}

func (p *pipeline) SetCollector(c *collector.Collector) *pipeline {
	p.collector = c

	return p
}

func (p *pipeline) ToInsert() *Insert {
	return newInsertFromPipeline(p)
}

func (p *pipeline) ToUpdate() *Update {
	return newUpdateFromPipeline(p)
}