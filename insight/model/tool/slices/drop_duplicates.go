package slices

type comparator interface {
	any | string | int
}

func SliceSliceDropDuplicates[T comparator](s [][]T, indexes []int) [][]T {
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
			_valueAny := any(_value)
			if _, ok := valuesFlag[index][_valueAny]; !ok {
				valuesFlag[index][_valueAny] = struct{}{}
				isDuplicates = false
			}
		}

		if !isDuplicates {
			newS = append(newS, item)
		}
	}

	return newS
}
