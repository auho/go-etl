package mode

import (
	"fmt"
	"strings"

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

func (m *InsertMode) GetTitle() string {
	return m.getModeTitle() + " " + m.insert.GetTitle()
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

type InsertMultiMode struct {
	Mode
	inserts      []means.InsertMeans
	insertFields []string
}

func NewInsertMultiMode(keys []string, insertFields []string, inserts []means.InsertMeans) *InsertMultiMode {
	m := &InsertMultiMode{}
	m.keys = keys
	m.inserts = inserts
	m.insertFields = insertFields

	return m
}

func (m *InsertMultiMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range m.inserts {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("%s{%s}", m.getModeTitle(), strings.Join(is, ","))
}

func (m *InsertMultiMode) GetFields() []string {
	return m.keys
}

func (m *InsertMultiMode) GetKeys() []string {
	return m.insertFields
}

func (m *InsertMultiMode) Do(item map[string]interface{}) [][]interface{} {
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

func (m *InsertMultiMode) Close() {
	for _, i := range m.inserts {
		i.Close()
	}
}
