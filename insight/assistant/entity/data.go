package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*Data)(nil)

type Data struct {
	base
	extra
	name   string
	idName string
}

func NewData(name string, idName string, db *simpledb.SimpleDB) *Data {
	d := &Data{}
	d.name = name
	d.idName = idName
	d.db = db
	d.extra = extra{
		base: d,
	}

	return d
}

func (d *Data) DB() *simpledb.SimpleDB {
	return d.db
}

func (d *Data) Name() string {
	return d.name
}

func (d *Data) IDName() string {
	return d.idName
}

func (d *Data) TableName() string {
	return fmt.Sprintf("%s_%s", NameData, d.name)
}

func (d *Data) WithCommand(fn func(command *schema.Command)) *Data {
	d.withCommand(fn)

	return d
}

func (d *Data) ToRows() *Rows {
	return NewRows(d.name, d.idName, d.db)
}

func (d *Data) ToRaw() *Raw {
	return NewRaw(d.name, d.db)
}
