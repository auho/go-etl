package mode

import (
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

// Transfer
// transfer data source to dst
type Transfer struct {
	Mode
	db          *goSimpleDb.SimpleDB
	fields      []string          // source data table filed name
	alias       map[string]string // alias map[data key name]alias key name
	fixedFields []string
	fixedData   map[string]any // fixed data map[data key name]value
}

// NewTransfer
// tableName source table name
func NewTransfer(db *goSimpleDb.SimpleDB, tableName string, alias map[string]string, fixedData map[string]any) *Transfer {
	m := &Transfer{}
	m.db = db
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
		m.fields, err = m.db.GetTableColumns(tableName)
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

func (m *Transfer) Do(item map[string]any) []map[string]any {
	for _, k := range m.fixedFields {
		item[k] = m.fixedData[k]
	}

	return []map[string]any{item}
}

func (m *Transfer) Close() {}
