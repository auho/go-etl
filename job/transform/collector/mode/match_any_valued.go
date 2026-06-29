package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// matchAnyValued iterates only keys with non-empty values (semantic B: by value).
// It skips empty-value keys, searches each remaining key's content one by one,
// and stops at the first match.
type matchAnyValued struct{}

var _ Mode = (*matchAnyValued)(nil)

func NewMatchAnyValued() Mode { return &matchAnyValued{} }

func (m *matchAnyValued) Prepare() error { return nil }
func (m *matchAnyValued) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	var st extract.Result
	for _, key := range valuedKeys(keys, keysValue) {
		st = e.Search([]string{keysValue[key]})
		if st.IsOK() {
			break
		}
	}
	return st, nil
}
