package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Moder = (*DataContentSpiltWords)(nil)

type DataContentSpiltWords struct {
	model
	data        *Data
	contentName string
}

func NewDataContentSpiltWords(data *Data, contentName string, db *simpleDb.SimpleDB) *DataContentSpiltWords {
	dc := &DataContentSpiltWords{}
	dc.data = data
	dc.contentName = contentName
	dc.db = db

	return dc
}

func (dc *DataContentSpiltWords) GetDB() *simpleDb.SimpleDB {
	return dc.db
}

func (dc *DataContentSpiltWords) GetName() string {
	return fmt.Sprintf("%s_%s", dc.data.GetName(), dc.contentName)
}

func (dc *DataContentSpiltWords) GetIdName() string {
	return "id"
}

func (dc *DataContentSpiltWords) GetData() *Data {
	return dc.data
}

func (dc *DataContentSpiltWords) GetContentName() string {
	return dc.contentName
}

func (dc *DataContentSpiltWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dc.data.GetName(), dc.contentName, NameSpiltWords)
}

func (dc *DataContentSpiltWords) WordName() string {
	return NameWord
}

func (dc *DataContentSpiltWords) WithCommand(fn func(command *tablestructure.Command)) *DataContentSpiltWords {
	dc.withCommand(fn)

	return dc
}
