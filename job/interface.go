package job

import (
	simpledb "github.com/auho/go-simple-db/v2"
)

type Source interface {
	GetIdName() string
	TableName() string
	GetDB() *simpledb.SimpleDB
}

type Target interface {
	GetIdName() string
	TableName() string
	GetDB() *simpledb.SimpleDB
}

type CleanResource interface {
	DataTarget() Target
	DeletedTarget() Target
	SourceTarget() Target
}
