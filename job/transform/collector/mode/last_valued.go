package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// lastValued selects the last key with a non-empty value (semantic B: by value).
// It filters out empty-value keys, then takes the last one.
type lastValued struct{}

var _ Mode = (*lastValued)(nil)

func NewLastValued() Mode { return &lastValued{} }

func (m *lastValued) Prepare() error { return nil }
func (m *lastValued) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Search(valuesByKeys(takeLast(valuedKeys(keys, keysValue), 1), keysValue)), nil
}
