package slices

func SliceToAny[T Any](s []T) []any {
	var newS []any
	for _, item := range s {
		newS = append(newS, item)
	}

	return newS
}

func SliceSliceToAny[T any](ss [][]T) [][]any {
	var newSS [][]any
	for _, s := range ss {
		var newS []any
		for _, item := range s {
			newS = append(newS, item)
		}

		newSS = append(newSS, newS)
	}

	return newSS
}
