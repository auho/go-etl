package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
)

// lastN selects the last n keys by position (semantic A: by key position).
// It searches the contents of keys[len-n:], regardless of whether values are empty.
type lastN struct{ n int }

var _ Mode = (*lastN)(nil)

func NewLastN(n int) Mode { return &lastN{n: n} }

func (m *lastN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *lastN) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Extract(valuesByKeys(takeLast(keys, m.n), keysValue)), nil
}
