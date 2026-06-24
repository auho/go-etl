package source

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*PlaceholderSource)(nil)

type PlaceholderSource struct {
	SourceBase
	baseCross
	basePlaceHolder
	items []map[string]any // []map[field][field value]
}

func NewPlaceholder(s SourceBase) *PlaceholderSource {
	return &PlaceholderSource{
		SourceBase: s,
	}
}

// WithItems
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

func (ps *PlaceholderSource) WithItems(items []map[string]any) *PlaceholderSource {
	ps.items = append(ps.items, items...)

	return ps
}

// WithItemsCross
// []map[string][]any => []map[field][][field value]
//
//	[]map[string][]any{
//		"one": []any{"a", "b"}
//		"two": []any{"c", "d"}
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
func (ps *PlaceholderSource) WithItemsCross(items map[string][]any) *PlaceholderSource {
	ps.WithItems(ps.expandItemsCross(items))

	return ps
}

func (ps *PlaceholderSource) Dataset() (*dataset.Dataset, error) {
	if len(ps.items) <= 0 {
		return nil, fmt.Errorf("PlaceholderSource source[%s] items len is error", ps.Name)
	}

	fields := ps.Table.GetSelectFields()
	sql := ps.Table.SQL()
	keys := ps.buildKeys(sql)

	itemsId, itemsSql := ps.buildPlaceholderItemsSqlSet(ps.SourceBase, sql, keys, ps.items)
	sets, err := ps.queryItemsSet(fields, itemsId, itemsSql)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet error; %w", err)
	}

	return &dataset.Dataset{
		Name:   ps.Name,
		Keys:   keys,
		Titles: fields,
		Sets:   sets,
	}, nil
}
