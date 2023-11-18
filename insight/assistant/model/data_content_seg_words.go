package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Moder = (*DataContentSegWords)(nil)

type DataContentSegWords struct {
	model
	data        assistant.Rowsor
	contentName string
}

func NewDataContentSegWords(data assistant.Rowsor, contentName string, db *simpleDb.SimpleDB) *DataContentSegWords {
	dc := &DataContentSegWords{}
	dc.data = data
	dc.contentName = contentName
	dc.db = db

	return dc
}

func (dc *DataContentSegWords) GetDB() *simpleDb.SimpleDB {
	return dc.db
}

func (dc *DataContentSegWords) GetName() string {
	return fmt.Sprintf("%s_%s", dc.data.GetName(), dc.contentName)
}

func (dc *DataContentSegWords) GetIdName() string {
	return "id"
}

func (dc *DataContentSegWords) GetData() assistant.Rowsor {
	return dc.data
}

func (dc *DataContentSegWords) GetContentName() string {
	return dc.contentName
}

func (dc *DataContentSegWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dc.data.GetName(), dc.contentName, NameSegWords)
}

func (dc *DataContentSegWords) WordName() string {
	return NameWord
}

func (dc *DataContentSegWords) FlagName() string {
	return NameFlag
}

func (dc *DataContentSegWords) NumName() string {
	return NameNum
}

func (dc *DataContentSegWords) WithCommand(fn func(command *tablestructure.Command)) *DataContentSegWords {
	dc.withCommand(fn)

	return dc
}
