package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*DataContent)(nil)

// DataContent represents a sub-table of a Data entity, associated with a specific content field.
// For example, a Data entity "articles" might have DataContent "articles_body" for the body content.
// The table name follows the pattern "data_<data_name>_<content_name>".
type DataContent struct {
	base        // embedded base for command hook and DB connection
	extra       // embedded extra for DDL/DML operations
	data        *Data  // parent Data entity
	contentName string // name of the content field (e.g., "body", "title")
}

// NewDataContent creates a new DataContent entity associated with the given Data and content field name.
func NewDataContent(data *Data, contentName string) *DataContent {
	d := &DataContent{}
	d.data = data
	d.contentName = contentName
	d.extra = extra{
		base: d,
	}

	return d
}

// DB returns the database connection from the parent Data entity.
func (d *DataContent) DB() *simpledb.SimpleDB {
	return d.data.DB()
}

// Name returns the combined name of the parent Data and the content field.
func (d *DataContent) Name() string {
	return fmt.Sprintf("%s_%s", d.data.name, d.contentName)
}

// IDName returns the primary key column name from the parent Data entity.
func (d *DataContent) IDName() string {
	return d.data.IDName()
}

// TableName returns the database table name, following the pattern "data_<data_name>_<content_name>".
func (d *DataContent) TableName() string {
	return fmt.Sprintf("%s_%s", NameData, d.Name())
}
