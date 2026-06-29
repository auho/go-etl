package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

type last struct{}

var _ Mode = (*last)(nil)

func NewLast() Mode { return &last{} }

func (m *last) Prepare() error { return nil }
func (m *last) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Search(valuesByKeys(takeLast(keys, 1), keysValue)), nil
}
