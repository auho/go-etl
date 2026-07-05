package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/task/extract"
)

// lastValuedN selects the last n keys with non-empty values (semantic B: by value).
// It filters out empty-value keys, then takes the last n of the remaining.
type lastValuedN struct{ n int }

var _ Mode = (*lastValuedN)(nil)

func NewLastValuedN(n int) Mode { return &lastValuedN{n: n} }

func (m *lastValuedN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *lastValuedN) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Extract(valuesByKeys(takeLast(valuedKeys(keys, keysValue), m.n), keysValue)), nil
}
