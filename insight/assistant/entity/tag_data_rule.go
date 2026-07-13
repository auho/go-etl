package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*TagDataRule)(nil)

// TagDataRule represents a single tagging result for a data and rule combination.
// Each TagDataRule corresponds to one rule applied to one data entity.
// The table name follows the pattern "tag_<data_name>_<rule_name>".
type TagDataRule struct {
	base                   // embedded base for command hook and DB connection
	extra                  // embedded extra for DDL/DML operations
	data  assistant.Entity // the data entity being tagged
	rule  assistant.Rule   // the rule applied to the data
}

// NewTagDataRule creates a new TagDataRule for the given data, rule, and database connection.
func NewTagDataRule(data assistant.Entity, rule assistant.Rule, db *simpledb.SimpleDB) *TagDataRule {
	t := &TagDataRule{
		base: base{db: db},
		data: data,
		rule: rule,
	}
	t.extra = extra{
		raw: t,
	}

	return t
}

// Data returns the data entity being tagged.
func (t *TagDataRule) Data() assistant.Entity {
	return t.data
}

// Rule returns the rule applied to the data.
func (t *TagDataRule) Rule() assistant.Rule {
	return t.rule
}

// Name returns the combined name of the data and the rule.
func (t *TagDataRule) Name() string {
	return fmt.Sprintf("%s_%s", t.data.Name(), t.rule.Name())
}

// IDName returns the primary key column name, which is always "id".
func (t *TagDataRule) IDName() string {
	return "id"
}

// TableName returns the database table name, following the pattern "tag_<data_name>_<rule_name>".
func (t *TagDataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.Name(), t.rule.Name())
}

// WithCommand sets the command hook and returns the TagDataRule instance for chaining.
func (t *TagDataRule) WithCommand(fn func(*schema.Command)) *TagDataRule {
	t.withCommand(fn)

	return t
}
