package assistant

import (
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
	simpledb "github.com/auho/go-simple-db/v3"
)

type Raw interface {
	GetDB() *simpledb.SimpleDB
	GetName() string
	TableName() string
	ExecCommand(*schema.Command) // exec command func
	DMLTable() *dml.Table
}

type Entity interface {
	Raw
	IDName() string
}

var _ Entity = Rule(nil)

type RuleConfig interface {
	AllowKeywordDuplicate() bool
}

type Rule interface {
	Entity
	NameLength() int
	Labels() map[string]int
	KeywordLength() int
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
