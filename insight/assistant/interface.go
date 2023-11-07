package assistant

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	"github.com/auho/go-simple-db/v2"
)

type Moder interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	GetIdName() string
	TableName() string
	ExecCommand(*tablestructure.Command)
}

type Rowsor interface {
	Moder
}

type Dataor interface {
	Moder
}

var _ Moder = Ruler(nil)

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
	ExecCommand(*tablestructure.Command)
}
