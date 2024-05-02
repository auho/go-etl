package explore

import (
	"maps"

	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.InsertModer = (*InsertStack)(nil)

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

func (is *InsertStack) Do(item map[string]any) []map[string]any {
	is.AddTotal(1)

	rets := make([]map[string]any, 0)
	for _, _i := range is.is {
		ret := _i.Do(item)
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

	return rets
}
