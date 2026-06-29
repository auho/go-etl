package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

type matchAnyN struct{ n int }

var _ collector.Mode = (*matchAnyN)(nil)

func NewMatchAnyN(n int) collector.Mode { return &matchAnyN{n: n} }

func (m *matchAnyN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *matchAnyN) Apply(source collector.Source, item map[string]any, e extract.Extractor) (extract.Result, error) {
	var st extract.Result
	for _, key := range takeFirst(source.Keys(), m.n) {
		v, err := source.Contents([]string{key}, item)
		if err != nil {
			return extract.Result{}, fmt.Errorf("contents: %w", err)
		}
		st = e.Search(v)
		if st.IsOK() {
			break
		}
	}
	return st, nil
}