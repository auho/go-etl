package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
)

var _ assistant.Ruler = (*Rule)(nil)

type DataRule struct {
	*Rule
	data assistant.Dataor
}

func NewDataRule(data assistant.Dataor, rule *Rule) *DataRule {
	dr := &DataRule{}
	dr.data = data
	dr.Rule = rule

	return dr
}

func (dr *DataRule) GetData() assistant.Dataor {
	return dr.data
}

func (dr *DataRule) GetRule() *Rule {
	return dr.Rule
}

func (dr *DataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameRule, dr.data.GetName(), dr.GetName())
}

func (dr *DataRule) ToOriginRule() assistant.Ruler {
	return dr.handlerOrigin()
}

func (dr *DataRule) ToAliasRule(alias map[string]string) *DataRule {
	_rule := dr.handlerOrigin()
	_rule.handlerAlias(alias)

	return _rule
}

func (dr *DataRule) ToItems(opts ...func(items *RuleItems)) *RuleItems {
	return NewRuleItems(dr, opts...)
}

func (dr *DataRule) handlerOrigin() *DataRule {
	return NewDataRule(dr.data, dr.Rule.handlerOrigin())
}
