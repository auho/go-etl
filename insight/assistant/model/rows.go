package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Rowsor = (*Rows)(nil)

type Rows struct {
	tableName string
	idName    string
	db        *simpleDb.SimpleDB
}

func NewRows(tableName string, db *simpleDb.SimpleDB) *Rows {
	return &Rows{
		tableName: tableName,
		db:        db,
	}
}

func (r *Rows) GetDB() *simpleDb.SimpleDB {
	return r.db
}

func (r *Rows) TableName() string {
	return r.tableName
}

func (r *Rows) DmlTable() *dml.Table {
	return dml.NewTable(r.TableName())
}
