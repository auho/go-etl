package transform

import (
	"fmt"
	"maps"
)

var _ InsertOperator = (*InsertStack)(nil)

// InsertStack
// 多个 insert 的所有结果 concat(上下拼接)
type InsertStack struct {
	baseInsert
}

func NewInsertStack(is ...*Insert) *InsertStack {
	return &InsertStack{
		baseInsert{
			name: "InsertStack",
			is:   is,
		},
	}
}

func (is *InsertStack) Apply(item map[string]any) ([]map[string]any, error) {
	is.AddTotal(1)

	rets := make([]map[string]any, 0)
	for _, _i := range is.is {
		ret, err := _i.Apply(item)
		if err != nil {
			return nil, fmt.Errorf("apply: %w", err)
		}
		if ret == nil {
			continue
		}

		for _, _r := range ret {
			_nr := make(map[string]any)
			maps.Copy(_nr, is.defaultValues)
			maps.Copy(_nr, _r)
			rets = append(rets, _nr)
		}
	}

	is.AddAmount(int64(len(rets)))

	return rets, nil
}
