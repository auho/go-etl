package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*DataContentSegWords)(nil)

// DataContentSegWords represents a table storing segmented words extracted from a content field.
// Segmentation is typically performed by NLP tools (e.g., jieba) that split text into meaningful terms.
// The table name follows the pattern "tag_<data_name>_<content_name>_seg_words".
type DataContentSegWords struct {
	base        // embedded base for command hook and DB connection
	extra       // embedded extra for DDL/DML operations
	data        assistant.Entity // parent data entity
	contentName string           // name of the content field being segmented
}

// NewDataContentSegWords creates a new DataContentSegWords entity.
func NewDataContentSegWords(data assistant.Entity, contentName string, db *simpledb.SimpleDB) *DataContentSegWords {
	dc := &DataContentSegWords{}
	dc.data = data
	dc.contentName = contentName
	dc.db = db
	dc.extra = extra{
		base: dc,
	}

	return dc
}

// DB returns the database connection.
func (dc *DataContentSegWords) DB() *simpledb.SimpleDB {
	return dc.db
}

// Name returns the combined name of the parent data and the content field.
func (dc *DataContentSegWords) Name() string {
	return fmt.Sprintf("%s_%s", dc.data.Name(), dc.contentName)
}

// IDName returns the primary key column name, which is always "id".
func (dc *DataContentSegWords) IDName() string {
	return "id"
}

// GetData returns the parent data entity.
func (dc *DataContentSegWords) GetData() assistant.Entity {
	return dc.data
}

// GetContentName returns the content field name.
func (dc *DataContentSegWords) GetContentName() string {
	return dc.contentName
}

// TableName returns the database table name, following the pattern "tag_<data_name>_<content_name>_seg_words".
func (dc *DataContentSegWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dc.data.Name(), dc.contentName, NameSegWords)
}

// WordName returns the column name for the word field.
func (dc *DataContentSegWords) WordName() string {
	return NameWord
}

// FlagName returns the column name for the flag/POS tag field.
func (dc *DataContentSegWords) FlagName() string {
	return NameFlag
}

// NumName returns the column name for the numeric count field.
func (dc *DataContentSegWords) NumName() string {
	return NameNum
}

// WithCommand sets the command hook and returns the DataContentSegWords instance for chaining.
func (dc *DataContentSegWords) WithCommand(fn func(command *schema.Command)) *DataContentSegWords {
	dc.withCommand(fn)

	return dc
}
