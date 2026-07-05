package transform

import (
	"fmt"
	"maps"
)

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

func (is *InsertSpread) Apply(item map[string]any) ([]map[string]any, error) {
	is.addTotal(1)

	_has := false
	ret := make(map[string]any, len(is.defaultValues))
	for _, _i := range is.is {
		res, err := _i.Apply(item)
		if err != nil {
			return nil, fmt.Errorf("apply: %w", err)
		}
		if res == nil {
			continue
		}

		_has = true
		maps.Copy(ret, res[0])
	}

	if _has {
		is.addAmount(1)

		_dv := maps.Clone(is.defaultValues)
		maps.Copy(_dv, ret)

		return []map[string]any{_dv}, nil
	} else {
		return nil, nil
	}
}
