package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

type lastN struct{ n int }

var _ collector.Mode = (*lastN)(nil)

func NewLastN(n int) collector.Mode { return &lastN{n: n} }

func (m *lastN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *lastN) Apply(source collector.Source, item map[string]any, e extract.Extractor) (extract.Result, error) {
	return searchContents(source, takeLast(source.Keys(), m.n), item, e)
}