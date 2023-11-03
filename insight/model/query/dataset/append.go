package dataset

import (
	"github.com/auho/go-etl/v2/insight/model/tool/slices"
)

type AppendMode struct {
	dataset *Dataset
}

func NewAppendMode(ds *Dataset) *AppendMode {
	return &AppendMode{dataset: ds}
}

func (af *AppendMode) Data() (*Data, error) {
	rows := [][]any{slices.SliceToAny(af.dataset.Titles)}

	for _, iRows := range af.dataset.ItemsSet {
		rows = append(rows, iRows...)
	}

	data := &Data{}
	data.Add(af.dataset.Name, rows)

	return data, nil
}
