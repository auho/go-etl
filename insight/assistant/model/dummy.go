package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Moder = (*Dummy)(nil)

type Dummy struct {
	model
	name   string
	idName string
}

func NewDummy(name string, idName string, db *simpleDb.SimpleDB) *Dummy {
	d := &Dummy{}
	d.name = name
	d.idName = idName
	d.db = db

	return d
}

func (d *Dummy) GetDB() *simpleDb.SimpleDB {
	return d.db
}

func (d *Dummy) GetName() string {
	return d.name
}

func (d *Dummy) GetIdName() string {
	return d.idName
}

func (d *Dummy) TableName() string {
	return d.name
}

func (d *Dummy) WithCommand(fn func(command *tablestructure.Command)) *Dummy {
	d.withCommand(fn)

	return d
}
