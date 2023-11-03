package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/dml"
	"github.com/auho/go-etl/v2/insight/model/query/dataset"
)

var _ Sourcer = (*TableSource)(nil)

type TableSource struct {
	Source
	Table dml.Tabler
}

func (ts *TableSource) Dataset() (*dataset.Dataset, error) {
	fields := ts.Table.GetSelectFields()
	itemsName := []string{ts.Name}

	itemsSet, err := ts.queryItemsSet(
		fields,
		itemsName,
		map[string]string{ts.Name: ts.Table.Sql()},
	)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet error; %w", err)
	}

	return &dataset.Dataset{
		Name:      ts.Name,
		Titles:    fields,
		ItemsName: itemsName,
		ItemsSet:  itemsSet,
	}, nil
}
