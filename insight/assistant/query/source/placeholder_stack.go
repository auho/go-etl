package source

import (
	"fmt"
	"maps"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

var _ Sourcer = (*PlaceholderStackSource)(nil)

/*
 c: 5, 6

 a: 1, 2
 b: 3, 4
=>
c: 5
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4

c: 6
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/

type PlaceholderStackSource struct {
	basePlaceHolder
	baseCross
	Source
	Table      dml.Tabler
	Categories []map[string]string // []map[field][field value]
	Stacks     []map[string]string // []map[field][field value]
}

func (pc *PlaceholderStackSource) Dataset() (*dataset.Dataset, error) {
	fields := pc.Table.GetSelectFields()

	var _sets []dataset.Set

	for _, _category := range pc.Categories {
		var _items []map[string]string

		for _, stack := range pc.Stacks {
			_item := make(map[string]string)
			maps.Copy(_item, _category)
			maps.Copy(_item, stack)

			_items = append(_items, _item)
		}

		_categoryPs := &PlaceholderSource{
			Source: pc.Source,
			Table:  pc.Table,
			Items:  _items,
		}

		_psDs, err := _categoryPs.Dataset()
		if err != nil {
			return nil, fmt.Errorf("dataset error; %w", err)
		}

		_sets = append(_sets, dataset.NewSetWithSets(pc.categoryToId(_category), _psDs.Sets))
	}

	return &dataset.Dataset{
		Name:   pc.Name,
		Titles: fields,
		Sets:   _sets,
	}, nil
}

func (pc *PlaceholderStackSource) categoryToId(category map[string]string) string {
	var keys []string
	for _, v := range category {
		keys = append(keys, v)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return pc.itemsNameToIdentification(keys)
}
