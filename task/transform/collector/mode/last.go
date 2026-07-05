package mode

import (
	"github.com/auho/go-etl/v3/task/extract"
)

// last selects the last key by position (semantic A: by key position).
// It searches the content of the last key, regardless of whether the value is empty.
type last struct{}

var _ Mode = (*last)(nil)

func NewLast() Mode { return &last{} }

func (m *last) Prepare() error { return nil }
func (m *last) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Extract(valuesByKeys(takeLast(keys, 1), keysValue)), nil
}
