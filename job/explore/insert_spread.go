package explore

import "maps"

// InsertSpread
// 取每个 insert 结果的第一条，进行 spread
type InsertSpread struct {
	baseInsert
}

func NewInsertSpread(is ...*Insert) *InsertSpread {
	return &InsertSpread{
		baseInsert{
			name: "InsertSpread",
			is:   is,
		},
	}
}

func (is *InsertSpread) Do(item map[string]any) []map[string]any {
	is.AddTotal(1)

	_has := false
	ret := make(map[string]any, len(is.defaultValues))
	for _, _i := range is.is {
		res := _i.Do(item)
		if res == nil {
			continue
		}

		_has = true
		maps.Copy(ret, res[0])
	}

	if _has {
		is.AddAmount(1)

		_dv := maps.Clone(is.defaultValues)
		maps.Copy(_dv, ret)

		return []map[string]any{_dv}
	} else {
		return nil
	}
}
