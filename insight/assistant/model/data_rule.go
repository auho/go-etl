package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

var _ assistant.Ruler = (*DataRule)(nil)

type DataRule struct {
	baseRule
	extra
	data assistant.Dataor
	rule *Rule
}

func NewDataRule(data assistant.Dataor, rule *Rule) *DataRule {
	dr := &DataRule{}
	dr.data = data
	dr.baseRule = rule.baseRule
	dr.rule = rule
	dr.extra = extra{
		model: dr,
	}

	return dr
}

func (dr *DataRule) handlerOrigin() *DataRule {
	return NewDataRule(dr.data, dr.rule.handlerOrigin())
}

func (dr *DataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameRule, dr.data.GetName(), dr.GetName())
}

func (dr *DataRule) ToOriginRule() assistant.Ruler {
	return dr.handlerOrigin()
}

func (dr *DataRule) ToItems(opts ...func(items *assistant.RuleItems)) *assistant.RuleItems {
	return assistant.NewRuleItems(dr, opts...)
}

func (dr *DataRule) WithCommand(fn func(command *tablestructure.Command)) *DataRule {
	dr.withCommand(fn)

	return dr
}

func (dr *DataRule) WithAllowKeywordDuplicate() *DataRule {
	dr.allowKeywordDuplicate = true

	return dr
}

func (dr *DataRule) ToAliasRule(alias map[string]string) *DataRule {
	_rule := dr.handlerOrigin()
	_rule.handlerAlias(alias)

	return _rule
}

func (dr *DataRule) GetData() assistant.Dataor {
	return dr.data
}

func (dr *DataRule) GetRule() *Rule {
	return dr.rule
}
