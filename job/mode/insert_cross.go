package mode

import (
	"maps"

	"github.com/auho/go-etl/v2/job/means"
)

var _ InsertModer = (*InsertCrossMode)(nil)

// InsertCrossMode
// cross means 交叉
//
// 1，2
// 3，4
// =>
// 1，3
// 1，4
// 2，3
// 2，4
type InsertCrossMode struct {
	insertHorizontalMode
}

func NewInsertCross(keys []string, ms ...means.InsertMeans) *InsertCrossMode {
	return &InsertCrossMode{
		insertHorizontalMode: newInsertHorizontal(keys, ms...),
	}
}

func (ic *InsertCrossMode) Do(item map[string]any) []map[string]any {
	ic.AddTotal(1)

	if item == nil {
		return nil
	}

	contents := ic.GetKeysContent(ic.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	var _allLabels [][]map[string]any
	for _, m := range ic.ms {
		mLabels := m.Insert(contents)
		if mLabels == nil {
			continue
		}

		_allLabels = append(_allLabels, mLabels)
	}

	var isStart = true
	var newItems []map[string]any
	var _tItems []map[string]any
	for _, _mLabels := range _allLabels {
		newItems = nil

		if isStart {
			isStart = false
			for _, _labels := range _mLabels {
				_nLabels := maps.Clone(ic.defaultValues)
				maps.Copy(_nLabels, _labels)
				newItems = append(newItems, _nLabels)
			}
		} else {
			for _, _tItem := range _tItems {
				for _, _resItem := range _mLabels {
					_newTItem := make(map[string]any)
					maps.Copy(_newTItem, _tItem)
					maps.Copy(_newTItem, _resItem)

					newItems = append(newItems, _newTItem)
				}
			}
		}

		_tItems = newItems
	}

	ic.AddAmount(int64(len(newItems)))

	return newItems
}
