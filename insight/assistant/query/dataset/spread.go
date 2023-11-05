package dataset

import (
	"github.com/auho/go-etl/v2/tool/slices"
)

type SpreadMode struct {
	dataset *Dataset
}

func NewSpreadMode(ds *Dataset) *SpreadMode {
	return &SpreadMode{dataset: ds}
}

func (af *SpreadMode) Data() (*Data, error) {
	data := &Data{}

	for _, itemName := range af.dataset.ItemsName {
		rows := [][]any{slices.SliceToAny(af.dataset.Titles)}

		data.Add(itemName, append(rows, af.dataset.ItemsSet[itemName]...))
	}

	return data, nil
}
