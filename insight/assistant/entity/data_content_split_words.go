package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*DataContentSplitWords)(nil)

// DataContentSplitWords represents a table storing split words from a content field.
// Unlike segmented words (NLP-based), split words are typically obtained by
// splitting text on delimiters (e.g., whitespace, punctuation).
// The table name follows the pattern "tag_<data_name>_<content_name>_split_words".
type DataContentSplitWords struct {
	base                         // embedded base for command hook and DB connection
	extra                        // embedded extra for DDL/DML operations
	data        assistant.Entity // parent data entity
	contentName string           // name of the content field being split
}

// NewDataContentSplitWords creates a new DataContentSplitWords entity.
func NewDataContentSplitWords(data assistant.Entity, contentName string, db *simpledb.SimpleDB) *DataContentSplitWords {
	dc := &DataContentSplitWords{
		base:        base{db: db},
		data:        data,
		contentName: contentName,
	}
	dc.extra = extra{
		raw: dc,
	}

	return dc
}

// Name returns the combined name of the parent data and the content field.
func (dc *DataContentSplitWords) Name() string {
	return fmt.Sprintf("%s_%s", dc.data.Name(), dc.contentName)
}

// IDName returns the primary key column name, which is always "id".
func (dc *DataContentSplitWords) IDName() string {
	return "id"
}

// TableName returns the database table name, following the pattern "tag_<data_name>_<content_name>_split_words".
func (dc *DataContentSplitWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dc.data.Name(), dc.contentName, NameSplitWords)
}

// GetData returns the parent data entity.
func (dc *DataContentSplitWords) GetData() assistant.Entity {
	return dc.data
}

// GetContentName returns the content field name.
func (dc *DataContentSplitWords) GetContentName() string {
	return dc.contentName
}

// WordName returns the column name for the word field.
func (dc *DataContentSplitWords) WordName() string {
	return NameWord
}

// WithCommand sets the command hook and returns the DataContentSplitWords instance for chaining.
func (dc *DataContentSplitWords) WithCommand(fn func(command *schema.Command)) *DataContentSplitWords {
	dc.withCommand(fn)

	return dc
}
