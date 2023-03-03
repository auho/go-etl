package mode

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/means"
)

type Insert struct {
	Mode
	insert means.InsertMeans
}

func NewInsert(keys []string, insert means.InsertMeans) *Insert {
	m := &Insert{}
	m.keys = keys
	m.insert = insert

	return m
}

func (m *Insert) GetTitle() string {
	return m.getModeTitle() + " " + m.insert.GetTitle()
}

func (m *Insert) GetFields() []string {
	return m.keys
}

func (m *Insert) GetKeys() []string {
	return m.insert.GetKeys()
}

func (m *Insert) Do(item map[string]interface{}) []map[string]interface{} {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	return m.insert.Insert(contents)
}

func (m *Insert) Close() {
	m.insert.Close()
}

type InsertMulti struct {
	Mode
	inserts      []means.InsertMeans
	insertFields []string
}

func NewInsertMulti(keys []string, insertFields []string, inserts []means.InsertMeans) *InsertMulti {
	m := &InsertMulti{}
	m.keys = keys
	m.inserts = inserts
	m.insertFields = insertFields

	return m
}

func (m *InsertMulti) GetTitle() string {
	is := make([]string, 0)
	for _, i := range m.inserts {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("%s{%s}", m.getModeTitle(), strings.Join(is, ","))
}

func (m *InsertMulti) GetFields() []string {
	return m.keys
}

func (m *InsertMulti) GetKeys() []string {
	return m.insertFields
}

func (m *InsertMulti) Do(item map[string]interface{}) []map[string]interface{} {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	items := make([]map[string]interface{}, 0)
	for _, i := range m.inserts {
		res := i.Insert(contents)
		if res == nil {
			continue
		}

		items = append(items, res...)
	}

	return items
}

func (m *InsertMulti) Close() {
	for _, i := range m.inserts {
		i.Close()
	}
}
