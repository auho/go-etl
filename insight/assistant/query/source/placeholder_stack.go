package source

import (
	"fmt"
	"maps"
	"sort"

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
	baseCross
	Source
	categories []map[string]any // []map[field][field value]
	stacks     []map[string]any // []map[field][field value]
}

func NewPlaceholderStack(s Source) *PlaceholderStackSource {
	return &PlaceholderStackSource{
		Source: s,
	}
}

func (pc *PlaceholderStackSource) WithCategories(categories []map[string]any) *PlaceholderStackSource {
	pc.categories = categories

	return pc
}

func (pc *PlaceholderStackSource) WithStacks(stacks []map[string]any) *PlaceholderStackSource {
	pc.stacks = stacks

	return pc
}

func (pc *PlaceholderStackSource) WithCategoriesCross(categories map[string][]any) *PlaceholderStackSource {
	return pc.WithCategories(pc.expandItemsCross(categories))
}

func (pc *PlaceholderStackSource) WithStacksCross(stacks map[string][]any) *PlaceholderStackSource {
	return pc.WithStacks(pc.expandItemsCross(stacks))
}

func (pc *PlaceholderStackSource) Dataset() (*dataset.Dataset, error) {
	if len(pc.categories) <= 0 {
		return nil, fmt.Errorf("PlaceholderStackSource[%s] categories len is error", pc.Name)
	}

	if len(pc.stacks) <= 0 {
		return nil, fmt.Errorf("PlaceholderStackSource[%s] stacks len is error", pc.Name)
	}

	fields := pc.Table.GetSelectFields()

	var _sets []dataset.Set

	for _, _category := range pc.categories {
		var _items []map[string]any

		for _, _stack := range pc.stacks {
			_item := make(map[string]any)
			maps.Copy(_item, _category)
			maps.Copy(_item, _stack)

			_items = append(_items, _item)
		}

		_categoryPs := NewPlaceholder(pc.Source).WithItems(_items)
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

func (pc *PlaceholderStackSource) categoryToId(category map[string]any) string {
	var keys []string
	for _, v := range category {
		keys = append(keys, fmt.Sprintf("%v", v))
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return pc.itemsNameToIdentification(keys)
}
