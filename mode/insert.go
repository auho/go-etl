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
