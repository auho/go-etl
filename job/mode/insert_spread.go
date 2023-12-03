package mode

import (
	"maps"

	"github.com/auho/go-etl/v2/job/means"
)

var _ InsertModer = (*InsertSpreadMode)(nil)

// InsertSpreadMode
// spread means
// 取每个 mean 结果的第一个，spread
type InsertSpreadMode struct {
	insertHorizontalMode
}

func NewInsertSpread(keys []string, ms ...means.InsertMeans) *InsertSpreadMode {
	return &InsertSpreadMode{
		insertHorizontalMode: newInsertHorizontal(keys, ms...),
	}
}

func (is *InsertSpreadMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := is.GetKeysContent(is.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	_has := false
	newItem := make(map[string]any, len(is.defaultValues))
	for _, m := range is.ms {
		res := m.Insert(contents)
		if res == nil {
			continue
		}

		_has = true

		maps.Copy(newItem, res[0])
	}

	if _has {
		_dv := maps.Clone(is.defaultValues)
		maps.Copy(_dv, newItem)

		return []map[string]any{_dv}
	} else {
		return nil
	}
}
