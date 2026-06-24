package transform

import (
	"maps"

	"github.com/auho/go-etl/v2/job/extract"
)

var _ InsertOperator = (*InsertSpread)(nil)

// InsertSpread
// spread means
// 取每个 mean 结果的第一个，spread
type InsertSpread struct {
	insertHorizontal
}

func NewInsertSpread(keys []string, ms ...extract.Inserter) *InsertSpread {
	return &InsertSpread{
		insertHorizontal: newInsertHorizontal(keys, ms...),
	}
}

func (is *InsertSpread) Do(item map[string]any) []map[string]any {
	is.AddTotal(1)

	if item == nil {
		return nil
	}

	contents := is.GetKeysContent(is.keys, item)
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
		is.AddAmount(1)

		_dv := maps.Clone(is.defaultValues)
		maps.Copy(_dv, newItem)

		return []map[string]any{_dv}
	} else {
		return nil
	}
}
