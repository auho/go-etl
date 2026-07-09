package entity

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Raw = (*Raw)(nil)

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

func (r *Raw) DB() *simpledb.SimpleDB {
	return r.db
}

func (r *Raw) Name() string {
	return r.name
}

func (r *Raw) TableName() string {
	return r.name
}
