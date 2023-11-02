package maps

func SliceMapStringAnyToSliceSliceAny(sm []map[string]any, keys []string) [][]any {
	var ss [][]any

	for _, row := range sm {
		var newRow []any
		for _, field := range keys {
			newRow = append(newRow, row[field])
		}

		ss = append(ss, newRow)
	}

	return ss
}
