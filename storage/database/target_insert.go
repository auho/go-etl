package database

func NewDbTargetInsertSliceSlice(config *DbTargetConfig, fields []string) *DbTargetSlice {
	if len(fields) <= 0 {
		panic("fields is error")
	}

	t := newDbTargetSlice(config)

	t.sliceFunc = func(d *DbTargetSlice, items [][]interface{}) error {
		_, err := d.db.BulkInsertFromSliceSlice(d.table, fields, items)

		return err
	}

	return t
}

func NewDbTargetInsertSliceMap(config *DbTargetConfig) *DbTargetMap {
	t := newDbTargetMap(config)

	t.mapFunc = func(d *DbTargetMap, items []map[string]interface{}) error {
		_, err := d.db.BulkInsertFromSliceMap(d.table, items)

		return err
	}

	return t
}
