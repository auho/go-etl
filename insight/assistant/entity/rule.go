package entity

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Rule = (*Rule)(nil)

// defaultStringLen is the default string length used for database column sizes
// when no explicit length is provided.
const defaultStringLen = 30

// Rule defines a rule for data analysis, consisting of a name, keyword configuration,
// and a set of labels. Rules are used to tag and analyze data based on keyword matching
// and label classification. The table name is prefixed with "rule_".
type Rule struct {
	baseRule // embedded rule base with name, labels, aliases, and keyword config
	extra    // embedded extra for DDL/DML operations
}

// NewRuleSimple creates a new Rule with default string lengths for all labels.
func NewRuleSimple(name string, labels []string, db *simpledb.SimpleDB) *Rule {
	_labels := make(map[string]int, len(labels))
	for _, label := range labels {
		_labels[label] = defaultStringLen
	}

	return NewRule(name, defaultStringLen, defaultStringLen, _labels, db)
}

// NewRule creates a new Rule with the given name, field lengths, labels, and DB connection.
func NewRule(name string, length, keywordLength int, labels map[string]int, db *simpledb.SimpleDB) *Rule {
	r := &Rule{}
	r.baseRule = newBaseRule(name, length, keywordLength, labels, db)
	r.extra = extra{
		raw: r,
	}

	return r
}

// handlerOrigin creates a new Rule with the original (non-aliased) name and labels.
// This is used as the base for alias operations.
func (r *Rule) handlerOrigin() *Rule {
	return NewRule(r.name, r.length, r.keywordLength, maps.Clone(r.labels), r.db)
}

// TableName returns the database table name, prefixed with "rule_".
// If an independent table name is set, it is used instead of the rule name.
func (r *Rule) TableName() string {
	_n := r.name
	if r.independentTableName != "" {
		_n = r.independentTableName
	}

	return fmt.Sprintf("%s_%s", NameRule, _n)
}

// ToOriginRule returns a new Rule with the original (non-aliased) name and labels.
func (r *Rule) ToOriginRule() assistant.Rule {
	return r.handlerOrigin()
}

// ToItems creates a RuleItems from this rule, optionally applying configuration options.
func (r *Rule) ToItems(opts ...func(items *assistant.RuleItems)) *assistant.RuleItems {
	return assistant.NewRuleItems(r, opts...)
}

// WithCommand sets the command hook and returns the Rule instance for chaining.
func (r *Rule) WithCommand(fn func(command *schema.Command)) *Rule {
	r.withCommand(fn)

	return r
}

// WithAllowKeywordDuplicate enables duplicate keyword support for this rule.
func (r *Rule) WithAllowKeywordDuplicate() *Rule {
	r.config.allowKeywordDuplicate = true

	return r
}

// WithTableName sets an independent table name (without the "rule_" prefix) for this rule.
// This overrides the default derived table name.
func (r *Rule) WithTableName(tableName string) *Rule {
	r.independentTableName = tableName

	return r
}

// ToAliasRule creates a new Rule with the given alias mapping applied.
// The alias maps original names to output names, allowing the same rule definition
// to produce differently named output tables.
func (r *Rule) ToAliasRule(alias map[string]string) *Rule {
	_rule := r.handlerOrigin()
	_rule.handlerAlias(alias)

	return _rule
}

// Clone creates a new Rule with the given name, preserving labels, aliases, and command hook.
func (r *Rule) Clone(name string) *Rule {
	return NewRule(name, r.length, r.keywordLength, r.labels, r.db).
		ToAliasRule(r.alias).
		WithCommand(r.commandFunc)
}

// CloneSuffix creates a new Rule by appending the given suffix to the current name.
func (r *Rule) CloneSuffix(suffix string) *Rule {
	return r.Clone(r.name + "_" + suffix)
}
