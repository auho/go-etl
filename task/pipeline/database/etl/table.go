package etl

import (
	simpledb "github.com/auho/go-simple-db/v3"
)

// Table exposes the table metadata and DB connection for ETL processing.
type Table interface {
	IDName() string
	TableName() string
	DB() *simpledb.SimpleDB
}

// CleanResource groups the three tables involved in a clean (filter) task:
// the source to read from, the data table for kept rows, and the deleted table
// for filtered rows.
type CleanResource interface {
	Data() Table
	Deleted() Table
	Source() Table
}
