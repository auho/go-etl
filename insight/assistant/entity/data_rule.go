package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

var _ assistant.Rule = (*DataRule)(nil)

// DataRule binds a Rule to a specific Data entity, creating a rule that operates
// on a particular data table. The table name follows the pattern "rule_<data_name>_<rule_name>".
// DataRule supports alias mapping and keyword duplicate configuration.
type DataRule struct {
	baseRule // embedded rule base with name, labels, aliases, and keyword config
	extra    // embedded extra for DDL/DML operations
	data     assistant.Entity // the data entity this rule is bound to
	rule     *Rule            // the underlying rule definition
}

// NewDataRule creates a new DataRule by binding the given rule to the given data entity.
func NewDataRule(data assistant.Entity, rule *Rule) *DataRule {
	dr := &DataRule{}
	dr.data = data
	dr.baseRule = rule.baseRule
	dr.rule = rule
	dr.extra = extra{
		base: dr,
	}

	return dr
}

// handlerOrigin creates a new DataRule with the original (non-aliased) names.
func (dr *DataRule) handlerOrigin() *DataRule {
	return NewDataRule(dr.data, dr.rule.handlerOrigin())
}

// TableName returns the database table name, following the pattern "rule_<data_name>_<rule_name>".
func (dr *DataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameRule, dr.data.Name(), dr.Name())
}

// ToOriginRule returns a new DataRule with the original (non-aliased) names.
func (dr *DataRule) ToOriginRule() assistant.Rule {
	return dr.handlerOrigin()
}

// ToItems creates a RuleItems from this rule, optionally applying configuration options.
func (dr *DataRule) ToItems(opts ...func(items *assistant.RuleItems)) *assistant.RuleItems {
	return assistant.NewRuleItems(dr, opts...)
}

// WithCommand sets the command hook and returns the DataRule instance for chaining.
func (dr *DataRule) WithCommand(fn func(command *schema.Command)) *DataRule {
	dr.withCommand(fn)

	return dr
}

// WithAllowKeywordDuplicate enables duplicate keyword support for this data rule.
func (dr *DataRule) WithAllowKeywordDuplicate() *DataRule {
	dr.config.allowKeywordDuplicate = true

	return dr
}

// ToAliasRule creates a new DataRule with the given alias mapping applied.
func (dr *DataRule) ToAliasRule(alias map[string]string) *DataRule {
	_rule := dr.handlerOrigin()
	_rule.handlerAlias(alias)

	return _rule
}

// GetData returns the data entity this rule is bound to.
func (dr *DataRule) GetData() assistant.Entity {
	return dr.data
}

// GetRule returns the underlying Rule definition.
func (dr *DataRule) GetRule() *Rule {
	return dr.rule
}
