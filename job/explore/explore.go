package explore

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v2/job/explore/collect"
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/mode"
)

type Explore struct {
	mode.Mode

	collect       collect.Collector
	search        search.Searcher
	defaultValues map[string]any
}

func NewExplorer(collect collect.Collector, search search.Searcher) *Explore {
	return &Explore{
		collect: collect,
		search:  search,
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
	e.defaultValues = e.search.GenExport().GetDefaultValues()

	return nil
}

func (e *Explore) State() []string {
	return []string{fmt.Sprintf("%s: %s", e.GetTitle(), e.GenCounter())}
}

func (e *Explore) Close() error { return nil }
