package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

var _ Sourcer = (*TableSource)(nil)

type TableSource struct {
	Source
	Table dml.Tabler
}

func (ts *TableSource) Dataset() (*dataset.Dataset, error) {
	fields := ts.Table.GetSelectFields()
	itemsName := []string{ts.Name}
	itemsSql := map[string]string{ts.Name: ts.Table.Sql()}

	sets, err := ts.queryItemsSet(
		fields,
		itemsName,
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
