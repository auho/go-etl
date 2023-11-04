package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/dml"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Dataor = (*Data)(nil)

type Data struct {
	name   string
	idName string
	db     *simpleDb.SimpleDB
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

func (d *Data) DmlTable() *dml.Table {
	return dml.NewTable(d.TableName())
}
