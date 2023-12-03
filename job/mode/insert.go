package mode

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
	"github.com/auho/go-etl/v2/tool/slices"
)

var _ InsertModer = (*InsertMode)(nil)
var _ InsertModer = (*InsertMultiMode)(nil)
var _ InsertModer = (*InsertSpreadMode)(nil)
var _ InsertModer = (*InsertCrossMode)(nil)

// InsertMode
// single means
type InsertMode struct {
	Mode
	means means.InsertMeans
}

func NewInsert(keys []string, means means.InsertMeans) *InsertMode {
	im := &InsertMode{}
	im.Keys = keys
	im.means = means

	return im
}

func (im *InsertMode) Prepare() error {
	if len(im.Keys) <= 0 {
		return fmt.Errorf("InsertMode Prepare kyes not exists error")
	}

	err := im.means.Prepare()
	if err != nil {
		return fmt.Errorf("InsertMode Prepare error; %w", err)
	}

	return nil
}

func (im *InsertMode) GetTitle() string {
	return im.GenTitle("InsertMode", im.means.GetTitle())
}

func (im *InsertMode) GetFields() []string {
	return im.Keys
}

func (im *InsertMode) GetKeys() []string {
	return im.means.GetKeys()
}

func (im *InsertMode) DefaultValues() map[string]any {
	return maps.Clone(im.means.DefaultValues())
}

func (im *InsertMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	return im.means.Insert(contents)
}

func (im *InsertMode) Close() error {
	err := im.means.Close()
	if err != nil {
		return fmt.Errorf("InsertMode close error; %w", err)
	}

	return nil
}

// InsertMultiMode
// multi means
// 多个 means append(上下拼接)，使用相同 column name
type InsertMultiMode struct {
	Mode
	meanses []means.InsertMeans

	insertKeys    []string
	defaultValues map[string]any
}

func NewInsertMulti(keys []string, meanses ...means.InsertMeans) *InsertMultiMode {
	im := &InsertMultiMode{}
	im.Keys = keys
	im.meanses = meanses

	return im
}

func (im *InsertMultiMode) Prepare() error {
	if len(im.Keys) <= 0 {
		return fmt.Errorf("InsertMultiMode Prepare kyes not exists error")
	}

	for _, m := range im.meanses {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("InsertMultiMode prepare error; %w", err)
		}
	}

	im.defaultValues = make(map[string]any)

	for _, m := range im.meanses {
		im.insertKeys = append(im.insertKeys, m.GetKeys()...)

		maps.Copy(im.defaultValues, m.DefaultValues())
	}

	im.insertKeys = slices.SliceDropDuplicates(im.insertKeys)

	return nil
}

func (im *InsertMultiMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range im.meanses {
		is = append(is, i.GetTitle())
	}

	return im.GenTitle("InsertMultiMode", strings.Join(is, ","))
}

func (im *InsertMultiMode) GetFields() []string {
	return im.Keys
}

func (im *InsertMultiMode) GetKeys() []string {
	return im.insertKeys
}

func (im *InsertMultiMode) DefaultValues() map[string]any {
	return maps.Clone(im.defaultValues)
}

func (im *InsertMultiMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	items := make([]map[string]any, 0)
	for _, m := range im.meanses {
		res := m.Insert(contents)
		if res == nil {
			continue
		}

		for _, _r := range res {
			_nr := make(map[string]any)
			maps.Copy(_nr, im.defaultValues)
			maps.Copy(_nr, _r)
			items = append(items, _nr)
		}
	}

	return items
}

func (im *InsertMultiMode) Close() error {
	for _, m := range im.meanses {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("InsertMultiMode close error; %w", err)
		}
	}

	return nil
}

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

func NewInsertCross(keys []string, meanses ...means.InsertMeans) *InsertCrossMode {
	return &InsertCrossMode{
		insertHorizontalMode: newInsertHorizontal(keys, meanses...),
	}
}

func (ic *InsertCrossMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := ic.GetKeysContent(ic.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	var _allLabels [][]map[string]any
	for _, m := range ic.meanses {
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

	return newItems
}

// InsertSpreadMode
// spread means
// 取每个 mean 结果的第一个，spread
type InsertSpreadMode struct {
	insertHorizontalMode
}

func NewInsertSpread(keys []string, meanses ...means.InsertMeans) *InsertSpreadMode {
	return &InsertSpreadMode{
		insertHorizontalMode: newInsertHorizontal(keys, meanses...),
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
	for _, m := range is.meanses {
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

// insertHorizontalMode
// 多个 means horizontal
type insertHorizontalMode struct {
	Mode
	meanses []means.InsertMeans

	insertKeys    []string
	defaultValues map[string]any
}

func newInsertHorizontal(keys []string, meanses ...means.InsertMeans) insertHorizontalMode {
	ih := insertHorizontalMode{}
	ih.Keys = keys
	ih.meanses = meanses

	return ih
}

func (ih *insertHorizontalMode) Prepare() error {
	if len(ih.Keys) <= 0 {
		return fmt.Errorf("insertHorizontalMode Prepare kyes not exists error")
	}

	for _, m := range ih.meanses {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare error; %w", err)
		}
	}

	ih.defaultValues = make(map[string]any)

	for _, m := range ih.meanses {
		ih.insertKeys = append(ih.insertKeys, m.GetKeys()...)

		maps.Copy(ih.defaultValues, m.DefaultValues())
	}

	return nil
}

func (ih *insertHorizontalMode) GetTitle() string {
	ss := make([]string, 0)
	for _, m := range ih.meanses {
		ss = append(ss, m.GetTitle())
	}

	return ih.GenTitle("insertHorizontalMode", strings.Join(ss, ","))
}

func (ih *insertHorizontalMode) GetFields() []string {
	return ih.Keys
}

func (ih *insertHorizontalMode) GetKeys() []string {
	return ih.insertKeys
}

func (ih *insertHorizontalMode) DefaultValues() map[string]any {
	return maps.Clone(ih.defaultValues)
}

func (ih *insertHorizontalMode) Close() error {
	for _, m := range ih.meanses {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("close error; %w", err)
		}
	}

	return nil
}
