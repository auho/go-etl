package dataset

import (
	"fmt"

	"github.com/auho/go-etl/v3/tool/slicex"
)

var _ Merger = (*SpreadMode)(nil)

// SpreadMode merges datasets by spreading each set into separate results.
type SpreadMode struct {
	dataset *Dataset
}

func NewSpreadMode(ds *Dataset) *SpreadMode {
	return &SpreadMode{dataset: ds}
}

func (sm *SpreadMode) Data() (*Result, error) {
	if sm.dataset == nil {
		return nil, fmt.Errorf("dataset is nil")
	}

	result := NewResult()

	titles := slicex.SliceToAny(sm.dataset.Titles)
	for _, set := range sm.dataset.Sets {
		result.addRowsWithTitles(set.Name, titles, set.Rows)
	}

	return result, nil
}

func (sm *SpreadMode) Name() string {
	return sm.dataset.Name
}

func (sm *SpreadMode) Sets() []Subset {
	return sm.dataset.Sets
}
