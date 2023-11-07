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
	data          *Data
	contentName   string
	contentLength int
}

func NewDataContentSegWords(data *Data, contentName string, contentLength int, db *simpleDb.SimpleDB) *DataContentSegWords {
	dc := &DataContentSegWords{}
	dc.data = data
	dc.contentName = contentName
	dc.contentLength = contentLength
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

func (dc *DataContentSegWords) CommandExec(command *tablestructure.Command) {
	dc.execCommand(command)
}

func (dc *DataContentSegWords) GetData() *Data {
	return dc.data
}

func (dc *DataContentSegWords) GetContentName() string {
	return dc.contentName
}

func (dc *DataContentSegWords) GetContentLength() int {
	return dc.contentLength
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
