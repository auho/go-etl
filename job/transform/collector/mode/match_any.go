package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

type matchAny struct{}

var _ collector.Mode = (*matchAny)(nil)

func NewMatchAny() collector.Mode { return &matchAny{} }

func (m *matchAny) Prepare() error { return nil }
func (m *matchAny) Apply(source collector.Source, item map[string]any, e extract.Extractor) (extract.Result, error) {
	var st extract.Result
	for _, key := range source.Keys() {
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