package dataset

import (
	"github.com/auho/go-etl/v2/tool/slices"
)

var _ Moder = (*SpreadMode)(nil)

// SpreadMode
// spread dataset
type SpreadMode struct {
	dataset *Dataset
}

func NewSpreadMode(ds *Dataset) *SpreadMode {
	return &SpreadMode{dataset: ds}
}

func (sm *SpreadMode) Data() (*Data, error) {
	data := &Data{}

	titles := slices.SliceToAny(sm.dataset.Titles)
	for _, set := range sm.dataset.Sets {
		var rows [][]any
		data.addRowsWithTitles(set.Name, titles, append(rows, set.Rows...))
	}

	return data, nil
}

func (sm *SpreadMode) Name() string {
	return sm.dataset.Name
}

func (sm *SpreadMode) Sets() []Set {
	return sm.dataset.Sets
}
