package source

import (
	"maps"
)

type baseCross struct {
}

// map[string][]any => []map[string]any
func (bc *baseCross) expandItemsCross(items map[string][]any) []map[string]any {
	/*
		a: 1, 2
		b: 3, 4

		step 1: expand first key
			a: 1
			a: 2

		step 2: cross with second key
			a: 1 b: 3
			a: 2 b: 3

		step 3: complete cross product
			a: 1 b: 3
			a: 2 b: 3
			a: 1 b: 4
			a: 2 b: 4
	*/

	var newItems []map[string]any
	var _tItems []map[string]any
	_isStart := true
	for key, values := range items {
		newItems = nil // clear to generate new combinations

		if _isStart { // first key
			_isStart = false
			for _, value := range values {
				newItems = append(newItems, map[string]any{key: value})
			}
		} else { // subsequent keys append
			for _, value := range values {
				for _, tItem := range _tItems { // all combinations from previous iteration
					_tItem := maps.Clone(tItem)
					_tItem[key] = value
					newItems = append(newItems, _tItem)
				}
			}
		}

		_tItems = newItems // keep current combinations for next iteration
	}

	return newItems
}
