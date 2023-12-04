package explore

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.InsertModer = (*Explore)(nil)

type Explore struct {
	mode.Mode

	collect  collect.Collector
	search   search.Searcher
	modeName string

	defaultValues map[string]any
}

func NewExplorer(collect collect.Collector, search search.Searcher, modeName string) *Explore {
	return &Explore{
		collect:  collect,
		search:   search,
		modeName: modeName,
	}
}

func (e *Explore) GetTitle() string {
	return e.GenTitle(e.collect.GetTitle(), e.search.GetTitle())
}

func (e *Explore) GetFields() []string {
	return e.collect.GetKeys()
}

func (e *Explore) GetKeys() []string {
	return nil
}

func (e *Explore) DefaultValues() map[string]any {
	return maps.Clone(e.defaultValues)
}

func (e *Explore) Prepare() error {
	e.defaultValues = e.search.DefaultTokenize().ToModer(e.modeName).DefaultValues()

	return nil
}

func (e *Explore) Do(item map[string]any) []map[string]any {
	e.AddTotal(1)

	_ticket := e.collect.Search(item, e.search.Search)
	if !_ticket.GetOk() {
		return nil
	}

	ret := _ticket.ToModer(e.modeName).ToTokenize()

	e.AddAmount(int64(len(ret)))

	return ret
}

func (e *Explore) State() []string {
	return []string{fmt.Sprintf("%s: %s", e.GetTitle(), e.GenCounter())}
}

func (e *Explore) Close() error { return nil }
