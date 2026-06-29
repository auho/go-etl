package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

type all struct{}

var _ collector.Mode = (*all)(nil)

func NewAll() collector.Mode { return &all{} }

func (m *all) Prepare() error { return nil }
func (m *all) Apply(source collector.Source, item map[string]any, e extract.Extractor) (extract.Result, error) {
	return searchContents(source, source.Keys(), item, e)
}