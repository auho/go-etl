package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// first selects the first key by position (semantic A: by key position).
// It searches the content of keys[0], regardless of whether the value is empty.
type first struct{}

var _ Mode = (*first)(nil)

func NewFirst() Mode { return &first{} }

func (m *first) Prepare() error { return nil }
func (m *first) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Extract(valuesByKeys(takeFirst(keys, 1), keysValue)), nil
}
