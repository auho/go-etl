package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

type first struct{}

var _ Mode = (*first)(nil)

func NewFirst() Mode { return &first{} }

func (m *first) Prepare() error { return nil }
func (m *first) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Search(valuesByKeys(takeFirst(keys, 1), keysValue)), nil
}
