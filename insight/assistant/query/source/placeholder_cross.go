package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

var _ Sourcer = (*PlaceholderCrossSource)(nil)

/*
 a: 1, 2
 b: 3, 4
=>
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/

type PlaceholderCrossSource struct {
	basePlaceHolder
	baseCross
	Source
	Table dml.Tabler
	Items []map[string][]string // []map[field][][field value]
}

func (pcs *PlaceholderCrossSource) Dataset() (*dataset.Dataset, error) {
	fields := pcs.Table.GetSelectFields()
	items := pcs.expandItems(pcs.Items)

	itemsName, itemsSql := pcs.buildPlaceholderItemsSqlSet(&pcs.Source, pcs.Table.Sql(), items)
	sets, err := pcs.queryItemsSet(fields, itemsName, itemsSql)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet error; %w", err)
	}

	return &dataset.Dataset{
		Name:   pcs.Name,
		Titles: fields,
		Sets:   sets,
	}, nil
}
