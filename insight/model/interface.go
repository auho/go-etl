package model

type Rowsor interface {
	TableName() string
}

type Ruler interface {
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
	GetName() string
	GetIdName() string
	TableName() string
}
