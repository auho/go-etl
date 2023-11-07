package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Rowsor = (*Rows)(nil)

type Rows struct {
	model
	tableName string
	idName    string
	db        *simpleDb.SimpleDB
}

func NewRows(tableName, idName string, db *simpleDb.SimpleDB) *Rows {
	return &Rows{
		tableName: tableName,
		idName:    idName,
		db:        db,
	}
}

func (r *Rows) GetDB() *simpleDb.SimpleDB {
	return r.db
}

func (r *Rows) GetName() string {
	return r.tableName
}

func (r *Rows) GetIdName() string {
	return r.idName
}

func (r *Rows) TableName() string {
	return r.tableName
}

func (r *Rows) CommandExec(command *tablestructure.Command) {
	r.execCommand(command)
}

func (r *Rows) WithCommand(fn func(command *tablestructure.Command)) *Rows {
	r.withCommand(fn)

	return r
}

func (r *Rows) DmlTable() *dml.Table {
	return dml.NewTable(r.TableName())
}
