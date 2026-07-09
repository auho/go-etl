package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*DataContentSplitWords)(nil)

type DataContentSplitWords struct {
	model
	extra
	data        assistant.Entity
	contentName string
}

func NewDataContentSplitWords(data assistant.Entity, contentName string, db *simpledb.SimpleDB) *DataContentSplitWords {
	dc := &DataContentSplitWords{}
	dc.data = data
	dc.contentName = contentName
	dc.db = db
	dc.extra = extra{
		model: dc,
	}

	return dc
}

func (dc *DataContentSplitWords) DB() *simpledb.SimpleDB {
	return dc.db
}

func (dc *DataContentSplitWords) Name() string {
	return fmt.Sprintf("%s_%s", dc.data.Name(), dc.contentName)
}

func (dc *DataContentSplitWords) IDName() string {
	return "id"
}

func (dc *DataContentSplitWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dc.data.Name(), dc.contentName, NameSplitWords)
}

func (dc *DataContentSplitWords) GetData() assistant.Entity {
	return dc.data
}

func (dc *DataContentSplitWords) GetContentName() string {
	return dc.contentName
}

func (dc *DataContentSplitWords) WordName() string {
	return NameWord
}

func (dc *DataContentSplitWords) WithCommand(fn func(command *schema.Command)) *DataContentSplitWords {
	dc.withCommand(fn)

	return dc
}
