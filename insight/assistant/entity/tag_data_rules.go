package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*TagDataRules)(nil)

// TagDataRules represents a tagging result for a data entity with multiple rules.
// Unlike TagDataRule (single rule), TagDataRules groups multiple rules under one name.
// The table name follows the pattern "tag_<data_name>_<name>".
type TagDataRules struct {
	base                   // embedded base for command hook and DB connection
	extra                  // embedded extra for DDL/DML operations
	name  string           // name of this tag group
	data  assistant.Entity // the data entity being tagged
	rules []assistant.Rule // the rules applied to the data
}

// NewTagDataRules creates a new TagDataRules for the given data, rules, and database connection.
func NewTagDataRules(name string, data assistant.Entity, rules []assistant.Rule, db *simpledb.SimpleDB) *TagDataRules {
	t := &TagDataRules{
		base:  base{db: db},
		name:  name,
		data:  data,
		rules: rules,
	}
	t.extra = extra{
		raw: t,
	}

	return t
}

// Data returns the data entity being tagged.
func (t *TagDataRules) Data() assistant.Entity {
	return t.data
}

// Rules returns the list of rules applied to the data.
func (t *TagDataRules) Rules() []assistant.Rule {
	return t.rules
}

// Name returns the combined name of the data and the tag group name.
func (t *TagDataRules) Name() string {
	return fmt.Sprintf("%s_%s", t.data.Name(), t.name)
}

// IDName returns the primary key column name, which is always "id".
func (t *TagDataRules) IDName() string {
	return "id"
}

// TableName returns the database table name, following the pattern "tag_<data_name>_<name>".
func (t *TagDataRules) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.Name(), t.name)
}

// WithCommand sets the command hook and returns the TagDataRules instance for chaining.
func (t *TagDataRules) WithCommand(fn func(*schema.Command)) *TagDataRules {
	t.withCommand(fn)

	return t
}

// Clone creates a new TagDataRules with the given name, preserving the data and rules.
func (t *TagDataRules) Clone(name string) *TagDataRules {
	return NewTagDataRules(name, t.data, t.rules, t.db)
}
