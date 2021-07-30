package database

func NewDbTargetInsertSliceSlice(config *DbTargetConfig, fields []string, prepareFuncs ...DbTargetSlicePrepareFunc) *DbTargetSlice {
	if len(fields) <= 0 {
		panic("fields is error")
	}

	t := newDbTargetSlice(config, prepareFuncs...)

	t.sliceFunc = func(d *DbTargetSlice, items [][]interface{}) error {
		_, err := d.db.BulkInsertFromSliceSlice(d.table, fields, items)

		return err
	}

	return t
}

func NewDbTargetInsertSliceMap(config *DbTargetConfig, prepareFuncs ...DbTargetMapPrepareFunc) *DbTargetMap {
	t := newDbTargetMap(config, prepareFuncs...)

	t.mapFunc = func(d *DbTargetMap, items []map[string]interface{}) error {
		_, err := d.db.BulkInsertFromSliceMap(d.table, items)

		return err
	}

	return t
}
