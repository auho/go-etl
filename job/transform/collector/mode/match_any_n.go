package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
)

// matchAnyN iterates the first n keys by position (semantic A: by key position).
// It searches each key's content one by one and stops at the first match.
// Empty-value keys are still searched.
type matchAnyN struct{ n int }

var _ Mode = (*matchAnyN)(nil)

func NewMatchAnyN(n int) Mode { return &matchAnyN{n: n} }

func (m *matchAnyN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *matchAnyN) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	var st extract.Result
	for _, key := range takeFirst(keys, m.n) {
		st = e.Extract([]string{keysValue[key]})
		if st.IsOK() {
			break
		}
	}
	return st, nil
}
