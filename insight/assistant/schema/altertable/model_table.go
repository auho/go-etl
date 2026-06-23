package altertable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v2"
)

type ModelTable struct {
	baseTable

	db *simpledb.SimpleDB
}

func NewModelTable(m assistant.Raw) *ModelTable {
	return &ModelTable{
		baseTable: newBaseTable(m.TableName()),
		db:        m.GetDB(),
	}
}

func (m *ModelTable) Build() error {
	return m.build(m.SQL(), m.db)
}

func (m *ModelTable) BuildAffixSql() ([]string, error) {
	_sql := m.SQL()
	return _sql, m.build(_sql, m.db)
}

func (m *ModelTable) BuildChange() error {
	return m.build(m.SqlForChange(), m.db)
}

func (m *ModelTable) BuildChangeAffixSql() ([]string, error) {
	_sql := m.SqlForChange()
	return _sql, m.build(_sql, m.db)
}

func (m *ModelTable) WithCommand(fn func(command *schema.Command)) *ModelTable {
	m.commandFunc = fn

	return m
}
