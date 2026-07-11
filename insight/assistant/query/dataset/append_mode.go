package dataset

import (
	"fmt"

	"github.com/auho/go-etl/v3/tool/slicex"
)

var _ Merger = (*AppendMode)(nil)

// AppendMode merges datasets by appending all rows under a single title.
type AppendMode struct {
	dataset *Dataset
}

func NewAppendMode(ds *Dataset) *AppendMode {
	return &AppendMode{dataset: ds}
}

func (am *AppendMode) Data() (*Result, error) {
	if am.dataset == nil {
		return nil, fmt.Errorf("dataset is nil")
	}

	var rows [][]any

	for _, set := range am.dataset.Sets {
		rows = append(rows, set.Rows...)
	}

	result := NewResult()
	result.addRowsWithTitles(am.dataset.Name, slicex.SliceToAny(am.dataset.Titles), rows)

	return result, nil
}

func (am *AppendMode) Name() string {
	return am.dataset.Name
}

func (am *AppendMode) Sets() []Subset {
	return am.dataset.Sets
}
