package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

var _ assistant.Rule = (*DataRule)(nil)

type DataRule struct {
	baseRule
	extra
	data assistant.Entity
	rule *Rule
}

func NewDataRule(data assistant.Entity, rule *Rule) *DataRule {
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
	return fmt.Sprintf("%s_%s_%s", NameRule, dr.data.Name(), dr.Name())
}

func (dr *DataRule) ToOriginRule() assistant.Rule {
	return dr.handlerOrigin()
}

func (dr *DataRule) ToItems(opts ...func(items *assistant.RuleItems)) *assistant.RuleItems {
	return assistant.NewRuleItems(dr, opts...)
}

func (dr *DataRule) WithCommand(fn func(command *schema.Command)) *DataRule {
	dr.withCommand(fn)

	return dr
}

func (dr *DataRule) WithAllowKeywordDuplicate() *DataRule {
	dr.config.allowKeywordDuplicate = true

	return dr
}

func (dr *DataRule) ToAliasRule(alias map[string]string) *DataRule {
	_rule := dr.handlerOrigin()
	_rule.handlerAlias(alias)

	return _rule
}

func (dr *DataRule) GetData() assistant.Entity {
	return dr.data
}

func (dr *DataRule) GetRule() *Rule {
	return dr.rule
}
