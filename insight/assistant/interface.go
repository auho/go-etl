package assistant

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	"github.com/auho/go-simple-db/v2"
)

type Rawer interface {
	GetDB() *go_simple_db.SimpleDB
	GetName() string
	TableName() string
	ExecCommand(*tablestructure.Command) // exec command func
	DmlTable() *dml.Table
}

type Moder interface {
	Rawer
	GetIdName() string
}

type Rowsor interface {
	Moder
}

type Dataor interface {
	Rowsor
}

var _ Moder = Ruler(nil)

type RuleConfigure interface {
	AllowKeywordDuplicate() bool
}

type Ruler interface {
	Moder
	GetNameLength() int
	GetLabels() map[string]int
	GetKeywordLength() int
	LabelsName() []string
	LabelsAlias() map[string]string
	LabelNumName() string
	KeywordName() string
	KeywordLenName() string
	KeywordNumName() string
	ToOriginRule() Ruler
	ToItems(opts ...func(items *RuleItems)) *RuleItems
	Config() RuleConfigure
}
