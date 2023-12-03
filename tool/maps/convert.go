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

func SliceMapToAny[KT keyEntity, VT valueEntity](sm []map[KT]VT) []map[KT]any {
	var nsm []map[KT]any
	for _, m := range sm {
		nsm = append(nsm, MapToAny(m))
	}

	return nsm
}

func MapToAny[KT keyEntity, VT valueEntity](m map[KT]VT) map[KT]any {
	nm := make(map[KT]any, len(m))
	for k, v := range m {
		nm[k] = v
	}

	return nm
}
