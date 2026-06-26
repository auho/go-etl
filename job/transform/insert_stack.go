package transform

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
	slices "github.com/auho/go-etl/v3/tool/slicex"
)

var _ InsertOperator = (*InsertStack)(nil)

// InsertStack
// stack inserter
// 多个 inserter append(上下拼接)，使用相同 column name
type InsertStack struct {
	base
	inserters []extract.Inserter

	insertKeys    []string
	defaultValues map[string]any
}

func NewInsertStack(keys []string, inserters ...extract.Inserter) *InsertStack {
	im := &InsertStack{}
	im.keys = keys
	im.inserters = inserters

	return im
}

func (im *InsertStack) Prepare() error {
	if len(im.keys) <= 0 {
		return fmt.Errorf("keys do not exist")
	}

	for _, m := range im.inserters {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare: %w", err)
		}
	}

	im.defaultValues = make(map[string]any)

	for _, m := range im.inserters {
		im.insertKeys = append(im.insertKeys, m.Keys()...)

		maps.Copy(im.defaultValues, m.DefaultValues())
	}

	im.insertKeys = slices.SliceDropDuplicates(im.insertKeys)

	return nil
}

func (im *InsertStack) Title() string {
	is := make([]string, 0)
	for _, i := range im.inserters {
		is = append(is, i.Title())
	}

	return im.GenTitle("InsertStack", strings.Join(is, ","))
}

func (im *InsertStack) GetFields() []string {
	return im.keys
}

func (im *InsertStack) Keys() []string {
	return im.insertKeys
}

func (im *InsertStack) DefaultValues() map[string]any {
	return maps.Clone(im.defaultValues)
}

func (im *InsertStack) Do(item map[string]any) []map[string]any {
	im.AddTotal(1)

	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.keys, item)
	if len(contents) <= 0 {
		return nil
	}

	items := make([]map[string]any, 0)
	for _, m := range im.inserters {
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

	im.AddAmount(int64(len(item)))

	return items
}

func (im *InsertStack) State() []string {
	return []string{fmt.Sprintf("%s: %s", im.Title(), im.GenCounter())}
}

func (im *InsertStack) Close() error {
	for _, m := range im.inserters {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("close: %w", err)
		}
	}

	return nil
}
