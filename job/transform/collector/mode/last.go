package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

type last struct{}

var _ collector.Mode = (*last)(nil)

func NewLast() collector.Mode { return &last{} }

func (m *last) Prepare() error { return nil }
func (m *last) Apply(source collector.Source, item map[string]any, e extract.Extractor) (extract.Result, error) {
	return searchContents(source, takeLast(source.Keys(), 1), item, e)
}