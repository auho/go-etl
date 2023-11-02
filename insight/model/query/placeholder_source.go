package query

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/insight/model/dml"
	"github.com/auho/go-etl/v2/insight/model/tool/maps"
	"github.com/auho/go-etl/v2/insight/model/tool/slices"
)

var _ sourcer = (*PlaceholderSource)(nil)

/*
	a: 1, 2
	b: 3, 4

	a: 1 b: 3
	a: 1 b: 3
	a: 2 b: 4
	a: 2 b: 4
*/

type PlaceholderSource struct {
	Source
	Table            dml.Tabler
	PlaceholderItems []map[string][]string
}

func (ps *PlaceholderSource) Rows() ([][]any, error) {
	return ps.rowsOfItems()
}

func (ps *PlaceholderSource) rowsOfItems() ([][]any, error) {
	var rows [][]any
	fields := ps.Table.GetSelectFields()
	rows = append(rows, slices.SliceToAny(fields))

	items := ps.expandItems()

	for _, item := range items {
		iRows, err := ps.rowsOfItem(item)
		if err != nil {
			return nil, err
		}

		rows = append(rows, iRows...)
	}

	return rows, nil
}

func (ps *PlaceholderSource) rowsOfItem(item map[string]string) ([][]any, error) {
	var rows [][]any

	sql := ps.Table.Sql()
	for key, value := range item {
		sql = strings.ReplaceAll(sql, fmt.Sprintf("##%s##", key), value)
	}

	sRows, err := ps.querySql(sql)
	if err != nil {
		return nil, fmt.Errorf("querySql error; %w", err)
	}

	rows = append(rows, sRows...)

	return rows, nil
}

func (ps *PlaceholderSource) expandItems() []map[string]string {
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
	for _, item := range ps.PlaceholderItems {
		items = nil // 清空，方便生成最新的组合
		for key, values := range item {
			if len(_tItems) == 0 { // 第一个 key
				for _, value := range values {
					items = append(items, map[string]string{key: value})
				}
			} else { // 之后的 key 追加
				for _, value := range values {
					for _, tItem := range _tItems { // 上一次大循环的所有组合
						tItem[key] = value
						items = append(items, tItem)
					}
				}
			}
		}

		_tItems = items // 保留当前的组合，方便后面进行追加新的组合
	}

	return items
}

func (ps *PlaceholderSource) querySql(sql string) ([][]any, error) {
	var rows []map[string]any
	err := ps.DB.Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("raw error; %w", err)
	}

	fields := ps.Table.GetSelectFields()

	return maps.SliceMapStringAnyToSliceSliceAny(rows, fields), nil
}
