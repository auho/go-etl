package dataset

import (
	slices "github.com/auho/go-etl/v3/tool/slicex"
)

var _ Mode = (*AppendMode)(nil)

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
