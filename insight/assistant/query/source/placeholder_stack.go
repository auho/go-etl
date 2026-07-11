package source

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*PlaceholderStack)(nil)

/*
 c: 5, 6

 a: 1, 2
 b: 3, 4
=>
c: 5
 a: 1 b: 3
 a: 1 b: 4
 a: 2 b: 3
 a: 2 b: 4

c: 6
 a: 1 b: 3
 a: 1 b: 4
 a: 2 b: 3
 a: 2 b: 4
*/

type PlaceholderStack struct {
	Base
	baseCross
	basePlaceholder
	categories []map[string]any // []map[field][field value]
	stacks     []map[string]any // []map[field][field value]
}

func NewPlaceholderStack(b Base) *PlaceholderStack {
	return &PlaceholderStack{
		Base: b,
	}
}

// AppendCategories appends categories to the existing categories list.
// Multiple calls accumulate categories.
//
//	[]map[string]any{
//		{"one": "a", "two": "c"},
//		{"one": "a", "two": "d"},
//		{"one": "b", "two": "c"},
//		{"one": "b", "two": "d"},
//	}
func (pss *PlaceholderStack) AppendCategories(categories []map[string]any) *PlaceholderStack {
	pss.categories = append(pss.categories, categories...)

	return pss
}

// SetCategories replaces all categories with the given categories.
func (pss *PlaceholderStack) SetCategories(categories []map[string]any) *PlaceholderStack {
	pss.categories = nil
	return pss.AppendCategories(categories)
}

// AppendStacks appends stacks to the existing stacks list.
// Multiple calls accumulate stacks.
//
//	[]map[string]any{
//		{"one": "a", "two": "c"},
//		{"one": "a", "two": "d"},
//		{"one": "b", "two": "c"},
//		{"one": "b", "two": "d"},
//	}
func (pss *PlaceholderStack) AppendStacks(stacks []map[string]any) *PlaceholderStack {
	pss.stacks = append(pss.stacks, stacks...)

	return pss
}

// SetStacks replaces all stacks with the given stacks.
func (pss *PlaceholderStack) SetStacks(stacks []map[string]any) *PlaceholderStack {
	pss.stacks = nil
	return pss.AppendStacks(stacks)
}

// AppendCategoriesCross expands categories via cross product and appends to the existing categories list.
// Multiple calls accumulate categories.
//
//	map[string][]any{
//		"one": []any{"a", "b"},
//		"two": []any{"c", "d"},
//	}
//
// =>
//
//	[]map[string]any{
//		{"one": "a", "two": "c"},
//		{"one": "a", "two": "d"},
//		{"one": "b", "two": "c"},
//		{"one": "b", "two": "d"},
//	}
func (pss *PlaceholderStack) AppendCategoriesCross(categories map[string][]any) *PlaceholderStack {
	return pss.AppendCategories(pss.expandItemsCross(categories))
}

// SetCategoriesCross replaces all categories with the cross-expanded categories.
func (pss *PlaceholderStack) SetCategoriesCross(categories map[string][]any) *PlaceholderStack {
	pss.categories = nil
	return pss.AppendCategoriesCross(categories)
}

// AppendStacksCross expands stacks via cross product and appends to the existing stacks list.
// Multiple calls accumulate stacks.
//
//	map[string][]any{
//		"one": []any{"a", "b"},
//		"two": []any{"c", "d"},
//	}
//
// =>
//
//	[]map[string]any{
//		{"one": "a", "two": "c"},
//		{"one": "a", "two": "d"},
//		{"one": "b", "two": "c"},
//		{"one": "b", "two": "d"},
//	}
func (pss *PlaceholderStack) AppendStacksCross(stacks map[string][]any) *PlaceholderStack {
	return pss.AppendStacks(pss.expandItemsCross(stacks))
}

// SetStacksCross replaces all stacks with the cross-expanded stacks.
func (pss *PlaceholderStack) SetStacksCross(stacks map[string][]any) *PlaceholderStack {
	pss.stacks = nil
	return pss.AppendStacksCross(stacks)
}

func (pss *PlaceholderStack) Dataset() (*dataset.Dataset, error) {
	if len(pss.categories) <= 0 {
		return nil, fmt.Errorf("source[%s] categories length is invalid", pss.Name)
	}

	if len(pss.stacks) <= 0 {
		return nil, fmt.Errorf("source[%s] stacks length is invalid", pss.Name)
	}

	fields := pss.Table.GetSelectFields()
	keys := pss.buildKeys(pss.Table.SQL())

	// remove duplicates
	_categoryIdMap := make(map[string]struct{})
	var _sets []dataset.Subset

	for _, _category := range pss.categories {
		var _items []map[string]any

		_categoryId := pss.categoryToID(_category, keys)
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

		_categoryPs := NewPlaceholder(pss.Base).AppendItems(_items)
		_categoryDataset, err := _categoryPs.Dataset()
		if err != nil {
			return nil, fmt.Errorf("dataset: %w", err)
		}

		_sets = append(_sets, dataset.NewSubsetFromSubsets(pss.categoryToID(_category, _categoryDataset.Keys), _categoryDataset.Sets))
	}

	return &dataset.Dataset{
		Name:   pss.Name,
		Titles: fields,
		Sets:   _sets,
	}, nil
}

func (pss *PlaceholderStack) categoryToID(category map[string]any, keys []string) string {
	var values []string

	for _, _k := range keys {
		if _cv, ok := category[_k]; ok {
			values = append(values, fmt.Sprintf("%v", _cv))
		}
	}

	return pss.itemValuesToIdentification(values)
}
