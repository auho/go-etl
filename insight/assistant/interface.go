package assistant

import (
	"github.com/auho/go-simple-db/v2"
)

type Rowsor interface {
	GetDB() *go_simple_db.SimpleDB
	GetIdName() string
	TableName() string
}

type Ruler interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	GetNameLength() int
	GetIdName() string
	GetLabels() map[string]int
	GetKeywordLength() int
	LabelsName() []string
	TableName() string
	KeywordName() string
	KeywordLenName() string
	KeywordNumName() string
}

type Dataor interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	GetIdName() string
	TableName() string
}
