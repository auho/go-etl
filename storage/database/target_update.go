package database

func NewDbTargetUpdateSliceMap(config *DbTargetConfig, idName string) *DbTargetMap {
	t := newDbTargetMap(config)

	t.mapFunc = func(d *DbTargetMap, items []map[string]interface{}) error {
		return d.db.BulkUpdateFromSliceMapById(d.table, idName, items)
	}

	return t
}
