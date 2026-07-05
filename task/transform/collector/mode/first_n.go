package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/task/extract"
)

// firstN selects the first n keys by position (semantic A: by key position).
// It searches the contents of keys[:n], regardless of whether values are empty.
type firstN struct{ n int }

var _ Mode = (*firstN)(nil)

func NewFirstN(n int) Mode { return &firstN{n: n} }

func (m *firstN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *firstN) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Extract(valuesByKeys(takeFirst(keys, m.n), keysValue)), nil
}
