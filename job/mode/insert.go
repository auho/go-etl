package mode

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
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
	im.keys = keys
	im.means = means

	return im
}

func (im *InsertMode) Prepare() error {
	err := im.means.Prepare()
	if err != nil {
		return fmt.Errorf("InsertMode Prepare error; %w", err)
	}

	return nil
}

func (im *InsertMode) GetTitle() string {
	return "InsertMode " + im.Mode.getTitle() + " " + im.means.GetTitle()
}

func (im *InsertMode) GetFields() []string {
	return im.keys
}

func (im *InsertMode) GetKeys() []string {
	return im.means.GetKeys()
}

func (im *InsertMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.keys, item)

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
// 多个 means merge，使用相同 key 名称，使用 第一个 means 的 config
type InsertMultiMode struct {
	Mode
	meanses []means.InsertMeans
}

func NewInsertMulti(keys []string, meanses ...means.InsertMeans) *InsertMultiMode {
	im := &InsertMultiMode{}
	im.keys = keys
	im.meanses = meanses

	return im
}

func (im *InsertMultiMode) Prepare() error {
	for _, m := range im.meanses {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("InsertMultiMode prepare error; %w", err)
		}
	}

	return nil
}

func (im *InsertMultiMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range im.meanses {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("InsertMultiMode %s{%s}", im.Mode.getTitle(), strings.Join(is, ","))
}

func (im *InsertMultiMode) GetFields() []string {
	return im.keys
}

func (im *InsertMultiMode) GetKeys() []string {
	return im.meanses[0].GetKeys()
}

func (im *InsertMultiMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.keys, item)

	items := make([]map[string]any, 0)
	for _, m := range im.meanses {
		res := m.Insert(contents)
		if res == nil {
			continue
		}

		items = append(items, res...)
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
// cross means
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

	contents := ic.GetKeysContent(ic.keys, item)

	var _allRes [][]map[string]any
	for _, m := range ic.meanses {
		res := m.Insert(contents)
		if res == nil {
			continue
		}

		_allRes = append(_allRes, res)
	}

	var isStart = true
	var newItems []map[string]any
	var _tItems []map[string]any
	for _, _res := range _allRes {
		newItems = nil

		if isStart {
			isStart = false
			newItems = append(newItems, _res...)
		} else {
			for _, _tItem := range _tItems {
				for _, _resItem := range _res {
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

	contents := is.GetKeysContent(is.keys, item)

	newItem := make(map[string]any, len(is.insertKeys))
	for _, m := range is.meanses {
		res := m.Insert(contents)
		if res == nil {
			continue
		}

		maps.Copy(newItem, res[0])
	}

	return []map[string]any{newItem}
}

// insertHorizontalMode
// 多个 means horizontal
type insertHorizontalMode struct {
	Mode
	meanses []means.InsertMeans

	insertKeys []string
}

func newInsertHorizontal(keys []string, meanses ...means.InsertMeans) insertHorizontalMode {
	ih := insertHorizontalMode{}
	ih.keys = keys
	ih.meanses = meanses

	return ih
}

func (ih *insertHorizontalMode) Prepare() error {
	for _, m := range ih.meanses {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare error; %w", err)
		}
	}

	for _, m := range ih.meanses {
		ih.insertKeys = append(ih.insertKeys, m.GetKeys()...)
	}

	return nil
}

func (ih *insertHorizontalMode) GetTitle() string {
	ss := make([]string, 0)
	for _, m := range ih.meanses {
		ss = append(ss, m.GetTitle())
	}

	return fmt.Sprintf("insertHorizontalMode %s{%s}", ih.Mode.getTitle(), strings.Join(ss, ","))
}

func (ih *insertHorizontalMode) GetFields() []string {
	return ih.keys
}

func (ih *insertHorizontalMode) GetKeys() []string {
	return ih.insertKeys
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
