package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Moder = (*DataContentSplitWords)(nil)

type DataContentSplitWords struct {
	model
	extra
	data        assistant.Rowsor
	contentName string
}

func NewDataContentSplitWords(data assistant.Rowsor, contentName string, db *simpledb.SimpleDB) *DataContentSplitWords {
	dc := &DataContentSplitWords{}
	dc.data = data
	dc.contentName = contentName
	dc.db = db
	dc.extra = extra{
		model: dc,
	}

	return dc
}

func (dc *DataContentSplitWords) GetDB() *simpledb.SimpleDB {
	return dc.db
}

func (dc *DataContentSplitWords) GetName() string {
	return fmt.Sprintf("%s_%s", dc.data.GetName(), dc.contentName)
}

func (dc *DataContentSplitWords) GetIDName() string {
	return "id"
}

func (dc *DataContentSplitWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dc.data.GetName(), dc.contentName, NameSplitWords)
}

func (dc *DataContentSplitWords) GetData() assistant.Rowsor {
	return dc.data
}

func (dc *DataContentSplitWords) GetContentName() string {
	return dc.contentName
}

func (dc *DataContentSplitWords) WordName() string {
	return NameWord
}

func (dc *DataContentSplitWords) WithCommand(fn func(command *tablestructure.Command)) *DataContentSplitWords {
	dc.withCommand(fn)

	return dc
}
