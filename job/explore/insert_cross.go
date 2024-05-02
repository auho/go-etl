package explore

import (
	"maps"

	"github.com/auho/go-etl/v2/job/mode"
)

var _ mode.InsertModer = (*InsertCross)(nil)

// InsertCross
// cross means 交叉
//
// 1，2
// 3，4
// =>
// 1，3
// 1，4
// 2，3
// 2，4
type InsertCross struct {
	baseInsert
}

func NewInsertCross(is ...*Insert) *InsertCross {
	return &InsertCross{
		baseInsert{
			name: "InsertCross",
			is:   is,
		},
	}
}

func (ic *InsertCross) Do(item map[string]any) []map[string]any {
	ic.AddTotal(1)

	var _allRet [][]map[string]any
	for _, m := range ic.is {
		_ret := m.Do(item)
		if _ret == nil {
			continue
		}

		_allRet = append(_allRet, _ret)
	}

	var isStart = true
	var rets []map[string]any
	var _tRets []map[string]any
	for _, _ret := range _allRet {
		rets = nil

		if isStart {
			isStart = false

			for _, _r := range _ret {
				_tr := maps.Clone(ic.defaultValues)
				maps.Copy(_tr, _r)
				rets = append(rets, _tr)
			}
		} else {
			for _, _tRet := range _tRets {
				for _, _r := range _ret {
					_tr := make(map[string]any)
					maps.Copy(_tr, _tRet)
					maps.Copy(_tr, _r)

					rets = append(rets, _tr)
				}
			}
		}

		_tRets = rets
	}

	ic.AddAmount(int64(len(rets)))

	return rets
}
