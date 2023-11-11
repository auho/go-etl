package source

import (
	"maps"
)

type baseCross struct {
}

// []map[string][]string => []map[field][][field value]
//
// []map[string]string => []map[field]field value
func (bc *baseCross) expandItems(items []map[string][]string) []map[string]string {
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

	var newItems []map[string]string

	var _tItems []map[string]string
	for _, item := range items {
		newItems = nil // 清空，方便生成最新的组合
		for key, values := range item {
			if len(_tItems) == 0 { // 第一个 key
				for _, value := range values {
					newItems = append(newItems, map[string]string{key: value})
				}
			} else { // 之后的 key 追加
				for _, value := range values {
					for _, tItem := range _tItems { // 上一次大循环的所有组合
						_tItem := maps.Clone(tItem)
						_tItem[key] = value
						newItems = append(newItems, _tItem)
					}
				}
			}
		}

		_tItems = newItems // 保留当前的组合，方便后面进行追加新的组合
	}

	return newItems
}
