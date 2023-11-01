package model

type Ruler interface {
	GetName() string
	GetLength() int
	GetLabels() map[string]int
	GetKeywordLength() int
	TableName() string
	KeywordLenName() string
	KeywordName() string
}

type Dataor interface {
	GetName() string
	GetIdName() string
	TableName() string
}
