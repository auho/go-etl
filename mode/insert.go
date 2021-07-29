package mode

import (
	"github.com/auho/go-etl/means"
)

type InsertMode struct {
	Mode
	insert means.InsertMeans
}

func NewInsertMode(keys []string, insert means.InsertMeans) *InsertMode {
	m := &InsertMode{}
	m.keys = keys
	m.insert = insert

	return m
}

func (m *InsertMode) GetFields() []string {
	return m.keys
}

func (m *InsertMode) GetKeys() []string {
	return m.insert.GetKeys()
}

func (m *InsertMode) Do(item map[string]interface{}) [][]interface{} {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	return m.insert.Insert(contents)
}

func (m *InsertMode) Close() {
	m.insert.Close()
}

type MultiInsertMode struct {
	Mode
	inserts      []means.InsertMeans
	insertFields []string
}

func NewMultiInsertMode(keys []string, insertFields []string, inserts []means.InsertMeans) *MultiInsertMode {
	m := &MultiInsertMode{}
	m.keys = keys
	m.inserts = inserts
	m.insertFields = insertFields

	return m
}

func (m *MultiInsertMode) GetFields() []string {
	return m.keys
}

func (m *MultiInsertMode) GetKey() []string {
	return m.insertFields
}

func (m *MultiInsertMode) Do(item map[string]interface{}) [][]interface{} {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	items := make([][]interface{}, 0)
	for _, i := range m.inserts {
		res := i.Insert(contents)
		if res == nil {
			continue
		}

		items = append(items, res...)
	}

	return items
}

func (m *MultiInsertMode) Close() {
	for _, i := range m.inserts {
		i.Close()
	}
}
