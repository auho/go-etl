package mode

import (
	go_simple_db "github.com/auho/go-simple-db/v2"
)

type Transfer struct {
	Mode
	Db          *go_simple_db.SimpleDB
	fields      []string
	alias       map[string]string
	fixedFields []string
	fixedData   map[string]interface{}
}

func NewTransfer(db *go_simple_db.SimpleDB, tableName string, alias map[string]string, fixedData map[string]interface{}) *Transfer {
	m := &Transfer{}
	m.Db = db
	m.keys = make([]string, 0)
	m.fields = make([]string, 0)
	m.alias = alias

	if len(m.alias) >= 0 {
		for k, v := range alias {
			m.fields = append(m.fields, k)
			m.keys = append(m.keys, v)
		}
	} else {
		var err error
		m.fields, err = db.GetTableColumns(tableName)
		if err != nil {
			panic(err)
		}

		for _, field := range m.fields {
			m.keys = append(m.keys, field)
		}
	}

	if len(fixedData) > 0 {
		for k := range fixedData {
			m.keys = append(m.keys, k)
			m.fixedFields = append(m.fixedFields, k)
		}

		m.fixedData = fixedData
	}

	return m
}

func (m *Transfer) GetTitle() string {
	return "Transfer"
}

func (m *Transfer) GetKeys() []string {
	return m.keys
}

func (m *Transfer) GetFields() []string {
	return m.fields
}

func (m *Transfer) Do(item map[string]interface{}) []map[string]interface{} {
	for _, k := range m.fixedFields {
		item[k] = m.fixedData[k]
	}

	return []map[string]interface{}{item}
}

func (m *Transfer) Close() {}
