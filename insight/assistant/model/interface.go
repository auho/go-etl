package model

import (
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Rowsor interface {
	GetDB() *simpleDb.SimpleDB
	TableName() string
}

type Ruler interface {
	GetDB() *simpleDb.SimpleDB
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
	GetDB() *simpleDb.SimpleDB
	GetName() string
	GetIdName() string
	TableName() string
}
