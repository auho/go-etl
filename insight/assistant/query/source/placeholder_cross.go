package source

import (
	"fmt"
	"maps"

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
	Source
	Table dml.Tabler
	Items []map[string][]string // []map[field][][field value]
}

func (pcs *PlaceholderCrossSource) Dataset() (*dataset.Dataset, error) {
	fields := pcs.Table.GetSelectFields()
	items := pcs.expandItems()

	itemsName, itemsSql := pcs.buildPlaceholderItemsSqlSet(pcs.Table.Sql(), items)
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

func (pcs *PlaceholderCrossSource) expandItems() []map[string]string {
	/*
		a: 1, 2
		b: 3, 4

		step -:
		a: 1
		a: 2

		step -:
		a: 1 b: 3
		a: 2 b: 3

		step -:
		a: 1 b: 3
		a: 2 b: 3
		a: 1 b: 4
		a: 2 b: 4
	*/

	var items []map[string]string

	var _tItems []map[string]string
	for _, item := range pcs.Items {
		items = nil // 清空，方便生成最新的组合
		for key, values := range item {
			if len(_tItems) == 0 { // 第一个 key
				for _, value := range values {
					items = append(items, map[string]string{key: value})
				}
			} else { // 之后的 key 追加
				for _, value := range values {
					for _, tItem := range _tItems { // 上一次大循环的所有组合
						_tItem := maps.Clone(tItem)
						_tItem[key] = value
						items = append(items, _tItem)
					}
				}
			}
		}

		_tItems = items // 保留当前的组合，方便后面进行追加新的组合
	}

	return items
}
