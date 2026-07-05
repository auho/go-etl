package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/task/extract"
)

// firstValuedN selects the first n keys with non-empty values (semantic B: by value).
// It filters out empty-value keys, then takes the first n of the remaining.
type firstValuedN struct{ n int }

var _ Mode = (*firstValuedN)(nil)

func NewFirstValuedN(n int) Mode { return &firstValuedN{n: n} }

func (m *firstValuedN) Prepare() error {
	if m.n <= 0 {
		return fmt.Errorf("n must be positive, got %d", m.n)
	}
	return nil
}
func (m *firstValuedN) Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error) {
	return e.Extract(valuesByKeys(takeFirst(valuedKeys(keys, keysValue), m.n), keysValue)), nil
}
