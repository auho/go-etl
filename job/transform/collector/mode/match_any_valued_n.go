package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
)

// matchAnyValuedN iterates the first n keys with non-empty values (semantic B: by value).
// It skips empty-value keys, takes the first n of the remaining, searches each
// one by one, and stops at the first match.
type matchAnyValuedN struct{ n int }

var _ Mode = (*matchAnyValuedN)(nil)

func NewMatchAnyValuedN(n int) Mode { return &matchAnyValuedN{n: n} }

func (m *matchAnyValuedN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *matchAnyValuedN) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	var st extract.Result
	for _, key := range takeFirst(valuedKeys(keys, keysValue), m.n) {
		st = e.Search([]string{keysValue[key]})
		if st.IsOK() {
			break
		}
	}
	return st, nil
}
