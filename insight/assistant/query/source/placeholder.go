package source

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*PlaceholderSource)(nil)

type PlaceholderSource struct {
	Base
	baseCross
	basePlaceholder
	items []map[string]any // []map[field][field value]
}

func NewPlaceholder(s Base) *PlaceholderSource {
	return &PlaceholderSource{
		Base: s,
	}
}

// AppendItems appends items to the existing items list.
// Multiple calls accumulate items.
//
//	[]map[string]any{
//		{"one": "a", "two": "c"},
//		{"one": "a", "two": "d"},
//		{"one": "b", "two": "c"},
//		{"one": "b", "two": "d"},
//	}

func (ps *PlaceholderSource) AppendItems(items []map[string]any) *PlaceholderSource {
	ps.items = append(ps.items, items...)

	return ps
}

// SetItems replaces all items with the given items.
func (ps *PlaceholderSource) SetItems(items []map[string]any) *PlaceholderSource {
	ps.items = nil
	return ps.AppendItems(items)
}

// AppendItemsCross expands items via cross product and appends to the existing items list.
// Multiple calls accumulate items.
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
func (ps *PlaceholderSource) AppendItemsCross(items map[string][]any) *PlaceholderSource {
	ps.AppendItems(ps.expandItemsCross(items))

	return ps
}

// SetItemsCross replaces all items with the cross-expanded items.
func (ps *PlaceholderSource) SetItemsCross(items map[string][]any) *PlaceholderSource {
	ps.items = nil
	return ps.AppendItemsCross(items)
}

func (ps *PlaceholderSource) Dataset() (*dataset.Dataset, error) {
	if len(ps.items) <= 0 {
		return nil, fmt.Errorf("source[%s] items length is invalid", ps.Name)
	}

	fields := ps.Table.GetSelectFields()
	sql := ps.Table.SQL()
	keys := ps.buildKeys(sql)

	itemsId, itemsSql := ps.buildPlaceholderItemsSqlSet(ps.Base, sql, keys, ps.items)
	sets, err := ps.queryItemsSet(fields, itemsId, itemsSql)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet: %w", err)
	}

	return &dataset.Dataset{
		Name:   ps.Name,
		Keys:   keys,
		Titles: fields,
		Sets:   sets,
	}, nil
}
