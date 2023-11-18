package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Dataor = (*Data)(nil)

type Data struct {
	model
	name   string
	idName string
}

func NewData(name string, idName string, db *simpleDb.SimpleDB) *Data {
	d := &Data{}
	d.name = name
	d.idName = idName
	d.db = db

	return d
}

func (d *Data) GetDB() *simpleDb.SimpleDB {
	return d.db
}

func (d *Data) GetName() string {
	return d.name
}

func (d *Data) GetIdName() string {
	return d.idName
}

func (d *Data) TableName() string {
	return fmt.Sprintf("%s_%s", NameData, d.name)
}

func (d *Data) WithCommand(fn func(command *tablestructure.Command)) *Data {
	d.withCommand(fn)

	return d
}

func (d *Data) DmlTable() *dml.Table {
	return dml.NewTable(d.TableName())
}

func (d *Data) ToRows() *Rows {
	return NewRows(d.name, d.idName, d.db)
}

func (d *Data) ToRaw() *Raw {
	return NewRaw(d.name, d.db)
}
