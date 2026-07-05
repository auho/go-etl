package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*DataContent)(nil)

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

func (d *DataContent) GetDB() *simpledb.SimpleDB {
	return d.data.GetDB()
}

func (d *DataContent) GetName() string {
	return fmt.Sprintf("%s_%s", d.data.name, d.contentName)
}

func (d *DataContent) IDName() string {
	return d.data.IDName()
}

func (d *DataContent) TableName() string {
	return fmt.Sprintf("%s_%s", NameData, d.GetName())
}
