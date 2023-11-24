package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Rawer = (*Raw)(nil)

type Raw struct {
	model
	name string
}

func NewRaw(name string, db *simpleDb.SimpleDB) *Raw {
	r := &Raw{}
	r.name = name
	r.db = db

	return r
}

func (r *Raw) GetDB() *simpleDb.SimpleDB {
	return r.db
}

func (r *Raw) GetName() string {
	return r.name
}

func (r *Raw) TableName() string {
	return r.name
}

func (r *Raw) DmlTable() *dml.Table {
	return dml.NewTable(r.TableName())
}
