package mode

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

var _ InsertModer = (*InsertMode)(nil)
var _ InsertModer = (*InsertMultiMode)(nil)

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
// 多个 means merge，使用相同 key 名称
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

func (m *InsertMultiMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range m.meanses {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("Insert multi %s{%s}", m.Mode.getTitle(), strings.Join(is, ","))
}

func (m *InsertMultiMode) GetFields() []string {
	return m.keys
}

func (m *InsertMultiMode) GetKeys() []string {
	return m.meanses[0].GetKeys()
}

func (m *InsertMultiMode) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	items := make([]map[string]any, 0)
	for _, i := range m.meanses {
		res := i.Insert(contents)
		if res == nil {
			continue
		}

		items = append(items, res...)
	}

	return items
}

func (m *InsertMultiMode) Close() error {
	for _, m := range m.meanses {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("InsertMultiMode close error; %w", err)
		}
	}

	return nil
}
