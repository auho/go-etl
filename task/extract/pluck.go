package extract

// PluckRows keeps only the specified keys from each row.
func PluckRows(rows []map[string]any, keys []string) []map[string]any {
	if len(keys) == 0 {
		return rows
	}

	result := make([]map[string]any, len(rows))
	for i, row := range rows {
		plucked := make(map[string]any, len(keys))
		for _, k := range keys {
			if v, ok := row[k]; ok {
				plucked[k] = v
			}
		}
		result[i] = plucked
	}
	return result
}
