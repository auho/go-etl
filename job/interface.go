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
	//GetGormDB() *gorm.DB
}

type CleanResource interface {
	Data() Target
	Deleted() Target
	Source() Target
}
