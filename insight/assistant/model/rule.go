package model

import (
	"maps"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Ruler = (*Rule)(nil)

const defaultStringLen = 30

type Rule struct {
	baseRule
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

	return r
}

func (r *Rule) ToOriginRule() assistant.Ruler {
	return r.handlerOrigin()
}

func (r *Rule) ToAliasRule(alias map[string]string) *Rule {
	_rule := r.handlerOrigin()
	_rule.handlerAlias(alias)

	return _rule
}

func (r *Rule) ToItems(opts ...func(items *RuleItems)) *RuleItems {
	return NewRuleItems(r, opts...)
}

func (r *Rule) DmlTable() *dml.Table {
	return dml.NewTable(r.TableName())
}

func (r *Rule) handlerOrigin() *Rule {
	return NewRule(r.name, r.length, r.keywordLength, maps.Clone(r.labels), r.db)
}
