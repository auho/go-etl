package mode

import (
	"github.com/auho/go-simple-db/simple"
)

type Transfer struct {
	Mode
	Db          simple.Driver
	fields      []string
	alias       map[string]string
	fixedValues []interface{}
}

func NewTransfer(db simple.Driver, tableName string, alias map[string]string, fixedData map[string]interface{}) *Transfer {
	m := &Transfer{}
	m.Db = db
	m.keys = make([]string, 0)
	m.fields = make([]string, 0)
	m.alias = alias
	m.fixedValues = make([]interface{}, 0)

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
		for k, v := range fixedData {
			m.keys = append(m.keys, k)
			m.fixedValues = append(m.fixedValues, v)
		}
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

func (m *Transfer) Do(item map[string]interface{}) [][]interface{} {
	result := make([]interface{}, len(m.fields), len(m.fields))
	for k, field := range m.fields {
		result[k] = item[field]
	}

	if len(m.fixedValues) > 0 {
		result = append(result, m.fixedValues...)
	}

	return [][]interface{}{result}
}

func (m *Transfer) Close() {

}
