package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Dataor = (*DataContent)(nil)

type DataContent struct {
	model
	extra
	data        *Data
	contentName string
}

func NewDataContent(data *Data, contentName string) *DataContent {
	d := &DataContent{}
	d.data = data
	d.contentName = contentName
	d.extra = extra{
		model: d,
	}

	return d
}

func (d *DataContent) GetDB() *simpleDb.SimpleDB {
	return d.data.GetDB()
}

func (d *DataContent) GetName() string {
	return fmt.Sprintf("%s_%s", d.data.name, d.contentName)
}

func (d *DataContent) GetIdName() string {
	return d.data.GetIdName()
}

func (d *DataContent) TableName() string {
	return fmt.Sprintf("%s_%s", NameData, d.GetName())
}

func (d *DataContent) DmlTable() *dml.Table {
	return dml.NewTable(d.TableName())
}
