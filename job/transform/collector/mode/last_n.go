package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
)

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
	return e.Search(valuesByKeys(takeLast(keys, m.n), keysValue)), nil
}
