package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// firstValued selects the first key with a non-empty value (semantic B: by value).
// It filters out empty-value keys, then takes the first one.
type firstValued struct{}

var _ Mode = (*firstValued)(nil)

func NewFirstValued() Mode { return &firstValued{} }

func (m *firstValued) Prepare() error { return nil }
func (m *firstValued) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Search(valuesByKeys(takeFirst(valuedKeys(keys, keysValue), 1), keysValue)), nil
}
