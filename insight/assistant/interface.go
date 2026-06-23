package assistant

import (
	"github.com/auho/go-etl/v2/insight/assistant/sqlbuilder/dml"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

type Raw interface {
	GetDB() *simpledb.SimpleDB
	GetName() string
	TableName() string
	ExecCommand(*tablestructure.Command) // exec command func
	DMLTable() *dml.Table
}

type Entity interface {
	Raw
	GetIDName() string
}

var _ Entity = Rule(nil)

type RuleConfig interface {
	AllowKeywordDuplicate() bool
}

type Rule interface {
	Entity
	GetNameLength() int
	GetLabels() map[string]int
	GetKeywordLength() int
	LabelsName() []string
	LabelsAlias() map[string]string
	LabelNumName() string
	TagsName() []string
	KeywordName() string
	KeywordLenName() string
	KeywordNumName() string
	KeywordAmountName() string
	ToOriginRule() Rule
	ToItems(opts ...func(items *RuleItems)) *RuleItems
	Config() RuleConfig
}
