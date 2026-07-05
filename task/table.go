package task

import (
	simpledb "github.com/auho/go-simple-db/v3"
)

type Table interface {
	IDName() string
	TableName() string
	GetDB() *simpledb.SimpleDB
}

type CleanResource interface {
	Data() Table
	Deleted() Table
	Source() Table
}
