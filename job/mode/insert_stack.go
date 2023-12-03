package mode

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
	"github.com/auho/go-etl/v2/tool/slices"
)

var _ InsertModer = (*InsertStackMode)(nil)

// InsertStackMode
// stack means
// 多个 means append(上下拼接)，使用相同 column name
type InsertStackMode struct {
	Mode
	ms []means.InsertMeans

	insertKeys    []string
	defaultValues map[string]any
}

func NewInsertStack(keys []string, ms ...means.InsertMeans) *InsertStackMode {
	im := &InsertStackMode{}
	im.Keys = keys
	im.ms = ms

	return im
}

func (im *InsertStackMode) Prepare() error {
	if len(im.Keys) <= 0 {
		return fmt.Errorf("InsertStackMode Prepare kyes not exists error")
	}

	for _, m := range im.ms {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("InsertStackMode prepare error; %w", err)
		}
	}

	im.defaultValues = make(map[string]any)

	for _, m := range im.ms {
		im.insertKeys = append(im.insertKeys, m.GetKeys()...)

		maps.Copy(im.defaultValues, m.DefaultValues())
	}

	im.insertKeys = slices.SliceDropDuplicates(im.insertKeys)

	return nil
}

func (im *InsertStackMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range im.ms {
		is = append(is, i.GetTitle())
	}

	return im.GenTitle("InsertStackMode", strings.Join(is, ","))
}

func (im *InsertStackMode) GetFields() []string {
	return im.Keys
}

func (im *InsertStackMode) GetKeys() []string {
	return im.insertKeys
}

func (im *InsertStackMode) DefaultValues() map[string]any {
	return maps.Clone(im.defaultValues)
}

func (im *InsertStackMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	items := make([]map[string]any, 0)
	for _, m := range im.ms {
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

func (im *InsertStackMode) Close() error {
	for _, m := range im.ms {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("InsertStackMode close error; %w", err)
		}
	}

	return nil
}
