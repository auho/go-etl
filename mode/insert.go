package mode

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/means"
)

// Insert
// single means
type Insert struct {
	Mode
	means means.InsertMeans
}

func NewInsert(keys []string, means means.InsertMeans) *Insert {
	m := &Insert{}
	m.keys = keys
	m.means = means

	return m
}

func (m *Insert) GetTitle() string {
	return "Insert " + m.Mode.getTitle() + " " + m.means.GetTitle()
}

func (m *Insert) GetFields() []string {
	return m.keys
}

func (m *Insert) GetKeys() []string {
	return m.means.GetKeys()
}

func (m *Insert) Do(item map[string]any) []map[string]any {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	return m.means.Insert(contents)
}

func (m *Insert) Close() {
	m.means.Close()
}

// InsertMulti
// multi means
// 多个 means merge，使用相同 key 名称
type InsertMulti struct {
	Mode
	meanses   []means.InsertMeans
	meansKeys []string
}

func NewInsertMulti(keys []string, meansKeys []string, meanses ...means.InsertMeans) *InsertMulti {
	m := &InsertMulti{}
	m.keys = keys
	m.meanses = meanses
	m.meansKeys = meansKeys

	return m
}

func (m *InsertMulti) GetTitle() string {
	is := make([]string, 0)
	for _, i := range m.meanses {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("Insert multi %s{%s}", m.Mode.getTitle(), strings.Join(is, ","))
}

func (m *InsertMulti) GetFields() []string {
	return m.keys
}

func (m *InsertMulti) GetKeys() []string {
	return m.meansKeys
}

func (m *InsertMulti) Do(item map[string]any) []map[string]any {
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

func (m *InsertMulti) Close() {
	for _, i := range m.meanses {
		i.Close()
	}
}
