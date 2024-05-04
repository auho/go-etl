package model

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Ruler = (*Rule)(nil)

const defaultStringLen = 30

type Rule struct {
	baseRule
	extra
}

func NewRuleSimple(name string, labels []string, db *simpleDb.SimpleDB) *Rule {
	_labels := make(map[string]int, len(labels))
	for _, label := range labels {
		_labels[label] = defaultStringLen
	}

	return NewRule(name, defaultStringLen, defaultStringLen, _labels, db)
}

func NewRule(name string, length, keywordLength int, labels map[string]int, db *simpleDb.SimpleDB) *Rule {
	r := &Rule{}
	r.baseRule = newBaseRule(name, length, keywordLength, labels, db)
	r.extra = extra{
		model: r,
	}

	return r
}

func (r *Rule) handlerOrigin() *Rule {
	return NewRule(r.name, r.length, r.keywordLength, maps.Clone(r.labels), r.db)
}

func (r *Rule) TableName() string {
	_n := r.name
	if r.independentTableName != "" {
		_n = r.independentTableName
	}

	return fmt.Sprintf("%s_%s", NameRule, _n)
}

func (r *Rule) ToOriginRule() assistant.Ruler {
	return r.handlerOrigin()
}

func (r *Rule) ToItems(opts ...func(items *assistant.RuleItems)) *assistant.RuleItems {
	return assistant.NewRuleItems(r, opts...)
}

func (r *Rule) WithCommand(fn func(command *tablestructure.Command)) *Rule {
	r.withCommand(fn)

	return r
}

func (r *Rule) WithAllowKeywordDuplicate() *Rule {
	r.config.allowKeywordDuplicate = true

	return r
}

// WithTableName
// 指定一个 table name 不包含 table prefix
func (r *Rule) WithTableName(tableName string) *Rule {
	r.independentTableName = tableName

	return r
}

// ToAliasRule
// db name 是 origin， output name 是 alias
func (r *Rule) ToAliasRule(alias map[string]string) *Rule {
	_rule := r.handlerOrigin()
	_rule.handlerAlias(alias)

	return _rule
}

// Clone
// only change name. copy labels, copy alias, copy command
func (r *Rule) Clone(name string) *Rule {
	return NewRule(name, r.length, r.keywordLength, r.labels, r.db).
		ToAliasRule(r.alias).
		WithCommand(r.commandFun)
}

func (r *Rule) CloneSuffix(suffix string) *Rule {
	return r.Clone(r.name + "_" + suffix)
}
