package slices

import (
	"fmt"
	"strings"
)

func SliceDropDuplicates[T comparatorEntity](s []T) []T {
	result := make([]T, 0)
	tempMap := make(map[any]bool, len(s))
	for _, e := range s {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}

	return result
}

// SliceSliceDropDuplicates
// indexes: index list
func SliceSliceDropDuplicates[T comparatorEntity](s [][]T, indexes []int) [][]T {
	var newS [][]T
	_sLen := len(s)
	valuesFlag := make(map[string]struct{}, _sLen)

	for _, item := range s {
		var _flags []string
		for _, _i := range indexes {
			_flags = append(_flags, fmt.Sprintf("%v", item[_i]))
		}

		_flag := strings.Join(_flags, "_")
		if _, ok := valuesFlag[_flag]; !ok {
			valuesFlag[_flag] = struct{}{}
			newS = append(newS, item)
		}
	}

	return newS
}
