package dataset

import (
	"github.com/auho/go-etl/v2/tool/slices"
)

var _ Moder = (*AppendMode)(nil)

// AppendMode
// append dataset
type AppendMode struct {
	dataset *Dataset
}

func NewAppendMode(ds *Dataset) *AppendMode {
	return &AppendMode{dataset: ds}
}

func (am *AppendMode) Data() (*Data, error) {
	var rows [][]any

	for _, set := range am.dataset.Sets {
		rows = append(rows, set.Rows...)
	}

	data := &Data{}
	data.addRowsWithTitles(am.dataset.Name, slices.SliceToAny(am.dataset.Titles), rows)

	return data, nil
}

func (am *AppendMode) Name() string {
	return am.dataset.Name
}

func (am *AppendMode) Sets() []Set {
	return am.dataset.Sets
}
