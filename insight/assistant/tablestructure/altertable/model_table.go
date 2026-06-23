package altertable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

type ModelTable struct {
	baseTable

	db *simpledb.SimpleDB
}

func NewModelTable(m assistant.Rawer) *ModelTable {
	return &ModelTable{
		baseTable: newBaseTable(m.TableName()),
		db:        m.GetDB(),
	}
}

func (m *ModelTable) Build() error {
	return m.build(m.Sql(), m.db)
}

func (m *ModelTable) BuildAffixSql() ([]string, error) {
	_sql := m.Sql()
	return _sql, m.build(_sql, m.db)
}

func (m *ModelTable) BuildChange() error {
	return m.build(m.SqlForChange(), m.db)
}

func (m *ModelTable) BuildChangeAffixSql() ([]string, error) {
	_sql := m.SqlForChange()
	return _sql, m.build(_sql, m.db)
}

func (m *ModelTable) WithCommand(fn func(command *tablestructure.Command)) *ModelTable {
	m.commandFunc = fn

	return m
}
