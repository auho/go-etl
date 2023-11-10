package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

var _ Sourcer = (*PlaceholderSource)(nil)

/*
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/

type PlaceholderSource struct {
	Source
	Table dml.Tabler
	Items []map[string]string // []map[field][field value]
}

func (ps *PlaceholderSource) Dataset() (*dataset.Dataset, error) {
	fields := ps.Table.GetSelectFields()

	itemsName, itemsSql := ps.buildPlaceholderItemsSqlSet(ps.Table.Sql(), ps.Items)
	sets, err := ps.queryItemsSet(fields, itemsName, itemsSql)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet error; %w", err)
	}

	return &dataset.Dataset{
		Name:   ps.Name,
		Titles: fields,
		Sets:   sets,
	}, nil
}
