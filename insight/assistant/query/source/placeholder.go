package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

var _ Sourcer = (*PlaceholderSource)(nil)

type PlaceholderSource struct {
	Source
	baseCross
	basePlaceHolder
	items []map[string]any // []map[field][field value]
}

func NewPlaceholder(s Source) *PlaceholderSource {
	return &PlaceholderSource{
		Source: s,
	}
}

// WithItems
// []map[string]string => []map[field][field value]
//
//	[]map[string]string{
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
	ps.items = items

	return ps
}

// WithItemsCross
// []map[string][]string => []map[field][][field value]
//
//	[]map[string][]string{
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
func (ps *PlaceholderSource) WithItemsCross(items map[string][]any) *PlaceholderSource {
	ps.items = ps.expandItemsCross(items)

	return ps
}

func (ps *PlaceholderSource) Dataset() (*dataset.Dataset, error) {
	if len(ps.items) <= 0 {
		return nil, fmt.Errorf("PlaceholderSource source[%s] items len is error", ps.Name)
	}

	fields := ps.Table.GetSelectFields()

	itemsName, itemsSql := ps.buildPlaceholderItemsSqlSet(&ps.Source, ps.Table.Sql(), ps.items)
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
