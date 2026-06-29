package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

type matchAny struct{}

var _ Mode = (*matchAny)(nil)

func NewMatchAny() Mode { return &matchAny{} }

func (m *matchAny) Prepare() error { return nil }
func (m *matchAny) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	var st extract.Result
	for _, key := range keys {
		st = e.Search([]string{keysValue[key]})
		if st.IsOK() {
			break
		}
	}
	return st, nil
}
