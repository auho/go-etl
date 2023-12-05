package altertable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

type ModeTable struct {
	baseTable

	db *simpledb.SimpleDB
}

func NewModelTable(m assistant.Rawer) *ModeTable {
	return &ModeTable{
		baseTable: newBaseTable(m.TableName()),
		db:        m.GetDB(),
	}
}

func (m *ModeTable) Build() error {
	return m.build(m.Sql(), m.db)
}

func (m *ModeTable) BuildAffixSql() ([]string, error) {
	_sql := m.Sql()
	return _sql, m.build(_sql, m.db)
}

func (m *ModeTable) BuildChange() error {
	return m.build(m.SqlForChange(), m.db)
}

func (m *ModeTable) BuildChangeAffixSql() ([]string, error) {
	_sql := m.SqlForChange()
	return _sql, m.build(_sql, m.db)
}

func (m *ModeTable) WithCommand(fn func(command *tablestructure.Command)) *ModeTable {
	m.commandFun = fn

	return m
}
