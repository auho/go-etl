package job

import (
	"github.com/auho/go-simple-db/v2"
)

type Source interface {
	GetIdName() string
	TableName() string
	GetDB() *go_simple_db.SimpleDB
}

type Target interface {
	GetIdName() string
	TableName() string
	GetDB() *go_simple_db.SimpleDB
}
