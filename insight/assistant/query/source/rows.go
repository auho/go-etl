package source

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*RowsSource)(nil)

// RowsSource
// general queries
type RowsSource struct {
	SourceBase
}

func NewRows(s SourceBase) *RowsSource {
	return &RowsSource{SourceBase: s}
}

func (ts *RowsSource) Dataset() (*dataset.Dataset, error) {
	fields := ts.Table.GetSelectFields()
	itemsId := []string{ts.Name}
	itemsSql := map[string]string{ts.Name: ts.Table.SQL()}

	sets, err := ts.queryItemsSet(
		fields,
		itemsId,
		itemsSql,
	)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet error; %w", err)
	}

	return &dataset.Dataset{
		Name:   ts.Name,
		Titles: fields,
		Sets:   sets,
	}, nil
}
