package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

type first struct{}

var _ collector.Mode = (*first)(nil)

func NewFirst() collector.Mode { return &first{} }

func (m *first) Prepare() error { return nil }
func (m *first) Apply(source collector.Source, item map[string]any, e extract.Extractor) (extract.Result, error) {
	return searchContents(source, takeFirst(source.Keys(), 1), item, e)
}