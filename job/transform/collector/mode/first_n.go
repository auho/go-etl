package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

type firstN struct{ n int }

var _ collector.Mode = (*firstN)(nil)

func NewFirstN(n int) collector.Mode { return &firstN{n: n} }

func (m *firstN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *firstN) Apply(source collector.Source, item map[string]any, e extract.Extractor) (extract.Result, error) {
	return searchContents(source, takeFirst(source.Keys(), m.n), item, e)
}