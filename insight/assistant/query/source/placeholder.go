package source

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*Placeholder)(nil)

type Placeholder struct {
	Base
	baseCross
	basePlaceholder
	items []map[string]any // []map[field][field value]
}

func NewPlaceholder(b Base) *Placeholder {
	return &Placeholder{
		Base: b,
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

func (ps *Placeholder) AppendItems(items []map[string]any) *Placeholder {
	ps.items = append(ps.items, items...)

	return ps
}

// SetItems replaces all items with the given items.
func (ps *Placeholder) SetItems(items []map[string]any) *Placeholder {
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
func (ps *Placeholder) AppendItemsCross(items map[string][]any) *Placeholder {
	ps.AppendItems(ps.expandItemsCross(items))

	return ps
}

// SetItemsCross replaces all items with the cross-expanded items.
func (ps *Placeholder) SetItemsCross(items map[string][]any) *Placeholder {
	ps.items = nil
	return ps.AppendItemsCross(items)
}

func (ps *Placeholder) Dataset() (*dataset.Dataset, error) {
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
