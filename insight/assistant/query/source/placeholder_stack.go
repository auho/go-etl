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
	basePlaceHolder
	Source
	categories []map[string]any // []map[field][field value]
	stacks     []map[string]any // []map[field][field value]
}

func NewPlaceholderStack(s Source) *PlaceholderStackSource {
	return &PlaceholderStackSource{
		Source: s,
	}
}

// WithCategories
// []map[string]any => []map[category][category value]
//
//	[]map[string]any{
//		{"one": "a", "two": "c"},
//		{"one": "a", "two": "d"},
//		{"one": "b", "two": "c"},
//		{"one": "b", "two": "d"},
//	}
/*
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/
func (pss *PlaceholderStackSource) WithCategories(categories []map[string]any) *PlaceholderStackSource {
	pss.categories = append(pss.categories, categories...)

	return pss
}

// WithStacks
// []map[string]any => []map[field][field value]
//
//	[]map[string]any{
//		{"one": "a", "two": "c"},
//		{"one": "a", "two": "d"},
//		{"one": "b", "two": "c"},
//		{"one": "b", "two": "d"},
//	}
/*
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/
func (pss *PlaceholderStackSource) WithStacks(stacks []map[string]any) *PlaceholderStackSource {
	pss.stacks = append(pss.stacks, stacks...)

	return pss
}

// WithCategoriesCross
// []map[string][]any => []map[category][][category value]
//
//	[]map[string][]any{
//		"one": []string{"a", "b"}
//		"two": []string{"c", "d"}
//	}
/*
 a: 1, 2
 b: 3, 4
=>
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/
func (pss *PlaceholderStackSource) WithCategoriesCross(categories map[string][]any) *PlaceholderStackSource {
	return pss.WithCategories(pss.expandItemsCross(categories))
}

// WithStacksCross
// []map[string][]any => []map[field][][field value]
//
//	[]map[string][]any{
//		"one": []string{"a", "b"}
//		"two": []string{"c", "d"}
//	}
/*
 a: 1, 2
 b: 3, 4
=>
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/
func (pss *PlaceholderStackSource) WithStacksCross(stacks map[string][]any) *PlaceholderStackSource {
	return pss.WithStacks(pss.expandItemsCross(stacks))
}

func (pss *PlaceholderStackSource) Dataset() (*dataset.Dataset, error) {
	if len(pss.categories) <= 0 {
		return nil, fmt.Errorf("PlaceholderStackSource[%s] categories len is error", pss.Name)
	}

	if len(pss.stacks) <= 0 {
		return nil, fmt.Errorf("PlaceholderStackSource[%s] stacks len is error", pss.Name)
	}

	fields := pss.Table.GetSelectFields()
	keys := pss.buildKeys(pss.Table.Sql())

	// remove duplicates
	_categoryIdMap := make(map[string]struct{})
	var _sets []dataset.Set

	for _, _category := range pss.categories {
		var _items []map[string]any

		_categoryId := pss.categoryToId(_category, keys)
		if _, ok := _categoryIdMap[_categoryId]; ok {
			continue
		}

		_categoryIdMap[_categoryId] = struct{}{}

		for _, _stack := range pss.stacks {
			_item := make(map[string]any)
			maps.Copy(_item, _category)
			maps.Copy(_item, _stack)

			_newItem := make(map[string]any, len(keys))
			for _, _k := range keys {
				_newItem[_k] = _item[_k]
			}

			_items = append(_items, _newItem)
		}

		_categoryPs := NewPlaceholder(pss.Source).WithItems(_items)
		_psDs, err := _categoryPs.Dataset()
		if err != nil {
			return nil, fmt.Errorf("dataset error; %w", err)
		}

		_sets = append(_sets, dataset.NewSetWithSets(pss.categoryToId(_category, _psDs.Keys), _psDs.Sets))
	}

	return &dataset.Dataset{
		Name:   pss.Name,
		Titles: fields,
		Sets:   _sets,
	}, nil
}

func (pss *PlaceholderStackSource) categoryToId(category map[string]any, keys []string) string {
	var values []string

	_keysMap := make(map[string]struct{}, len(keys))
	for _, _k := range keys {
		_keysMap[_k] = struct{}{}
	}

	for _ck, _cv := range category {
		if _, ok := _keysMap[_ck]; ok {
			values = append(values, fmt.Sprintf("%v", _cv))
		} else {
			//panic(fmt.Sprintf("categoryToId category[%s] value not found", _ck))
		}
	}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	return pss.itemValuesToIdentification(values)
}
