package assistant

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	"github.com/auho/go-simple-db/v2"
)

type Rawer interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	TableName() string
	ExecCommand(*tablestructure.Command)
}

type Moder interface {
	Rawer
	GetIdName() string
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
	LabelsAlias() map[string]string
	TableName() string
	KeywordName() string
	KeywordLenName() string
	KeywordNumName() string
	ExecCommand(*tablestructure.Command)
	ToOriginRule() Ruler
}
