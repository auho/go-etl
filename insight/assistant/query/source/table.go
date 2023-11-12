package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

var _ Sourcer = (*TableSource)(nil)

// TableSource
// general queries
type TableSource struct {
	Source
}

func NewTable(s Source) *TableSource {
	return &TableSource{Source: s}
}

func (ts *TableSource) Dataset() (*dataset.Dataset, error) {
	fields := ts.Table.GetSelectFields()
	itemsId := []string{ts.Name}
	itemsSql := map[string]string{ts.Name: ts.Table.Sql()}

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
