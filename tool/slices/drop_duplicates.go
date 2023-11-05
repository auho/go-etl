package slices

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
	valuesFlag := make(map[int]map[any]struct{}, len(indexes))

	for _, index := range indexes {
		valuesFlag[index] = make(map[any]struct{}, _sLen)
	}

	//var _valueAny any
	for _, item := range s {
		isDuplicates := true
		for _, index := range indexes {
			if !isDuplicates {
				break
			}

			_value := item[index]
			if _, ok := valuesFlag[index][_value]; !ok {
				valuesFlag[index][_value] = struct{}{}
				isDuplicates = false
			}
		}

		if !isDuplicates {
			newS = append(newS, item)
		}
	}

	return newS
}
