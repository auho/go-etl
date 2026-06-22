package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	simpledb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Rawer = (*Raw)(nil)

type Raw struct {
	model
	extra
	name string
}

func NewRaw(name string, db *simpledb.SimpleDB) *Raw {
	r := &Raw{}
	r.name = name
	r.db = db
	r.extra = extra{
		model: r,
	}

	return r
}

func (r *Raw) GetDB() *simpledb.SimpleDB {
	return r.db
}

func (r *Raw) GetName() string {
	return r.name
}

func (r *Raw) TableName() string {
	return r.name
}
