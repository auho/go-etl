package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// allValued selects only keys with non-empty values (semantic B: by value).
// It filters out empty-value keys, then searches all remaining contents together.
type allValued struct{}

var _ Mode = (*allValued)(nil)

func NewAllValued() Mode { return &allValued{} }

func (m *allValued) Prepare() error { return nil }
func (m *allValued) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Extract(valuesByKeys(valuedKeys(keys, keysValue), keysValue)), nil
}
