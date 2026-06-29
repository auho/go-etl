package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

type all struct{}

var _ Mode = (*all)(nil)

func NewAll() Mode { return &all{} }

func (m *all) Prepare() error { return nil }
func (m *all) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Search(valuesByKeys(keys, keysValue)), nil
}
