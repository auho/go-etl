package entity

import (
	"fmt"
	"sort"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.RuleConfig = (*RuleConfig)(nil)

// RuleConfig holds configuration options for rule behavior.
type RuleConfig struct {
	allowKeywordDuplicate bool // whether to allow duplicate keywords in analysis results
}

// AllowKeywordDuplicate returns whether duplicate keywords are allowed.
func (rc RuleConfig) AllowKeywordDuplicate() bool {
	return rc.allowKeywordDuplicate
}

// baseRule is the foundational embedded struct for rule entities (Rule, DataRule).
// It manages the rule's name, labels, keyword configuration, and alias mapping.
// Aliases allow the same rule to produce differently named output tables
// without duplicating the rule definition.
type baseRule struct {
	base
	config        RuleConfig
	name          string          // origin name of the rule
	length        int             // origin name string length (for database column size)
	keywordLength int             // keyword string length (for database column size)
	labels        map[string]int  // map[label]label length (for database column sizes)

	alias             map[string]string // map[origin name][alias name]
	nameAlias         string            // alias of the rule name (used as the effective name)
	labelsAlias       map[string]string // map[label]label alias
	labelsAliasLength map[string]int    // map[label alias]label alias length (for alias labels)

	independentTableName string // optional independent table name (overrides the default derived name)
}

// newBaseRule creates a new baseRule with the given name, field lengths, labels, and DB connection.
// Default string lengths are applied if the provided values are zero or negative.
func newBaseRule(name string, length, keywordLength int, labels map[string]int, db *simpledb.SimpleDB) baseRule {
	br := baseRule{}
	br.name = name
	br.length = length
	br.keywordLength = keywordLength
	br.labels = labels
	br.db = db

	br.handlerAlias(nil)

	if br.length <= 0 {
		br.length = defaultStringLen
	}

	if br.keywordLength <= 0 {
		br.keywordLength = defaultStringLen
	}

	return br
}

// DB returns the database connection.
func (br *baseRule) DB() *simpledb.SimpleDB {
	return br.db
}

// Name returns the effective name (alias if set, otherwise the original name).
func (br *baseRule) Name() string {
	return br.nameAlias
}

// NameLength returns the string length for the name field in the database.
func (br *baseRule) NameLength() int {
	return br.length
}

// IDName returns the primary key column name, which is always "id".
func (br *baseRule) IDName() string {
	return "id"
}

// KeywordLength returns the string length for the keyword field in the database.
func (br *baseRule) KeywordLength() int {
	return br.keywordLength
}

// Labels returns the effective label map (alias lengths if aliases are set).
func (br *baseRule) Labels() map[string]int {
	return br.labelsAliasLength
}

// TagsName returns a list of all tag names, including the rule name and all label names.
func (br *baseRule) TagsName() []string {
	var tagsName []string
	tagsName = append(tagsName, br.Name())
	tagsName = append(tagsName, br.LabelsName()...)

	return tagsName
}

// LabelsName returns a sorted list of all label names.
func (br *baseRule) LabelsName() []string {
	var labels []string
	for label, _ := range br.Labels() {
		labels = append(labels, label)
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	return labels
}

// LabelsAlias returns the label alias mapping (original name -> alias).
func (br *baseRule) LabelsAlias() map[string]string {
	return br.labelsAlias
}

// LabelNumName returns the column name for the label count field.
func (br *baseRule) LabelNumName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameLabelNum)
}

// KeywordName returns the column name for the keyword field.
func (br *baseRule) KeywordName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeyword)
}

// KeywordLenName returns the column name for the keyword length field.
func (br *baseRule) KeywordLenName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeywordLen)
}

// KeywordNumName returns the column name for the keyword count field.
func (br *baseRule) KeywordNumName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeywordNum)
}

// KeywordAmountName returns the column name for the keyword total amount field.
func (br *baseRule) KeywordAmountName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeywordAmount)
}

// Config returns the rule configuration.
func (br *baseRule) Config() assistant.RuleConfig {
	return br.config
}

// WithCommand sets the command hook and returns the baseRule for chaining.
func (br *baseRule) WithCommand(fn func(command *schema.Command)) *baseRule {
	br.withCommand(fn)

	return br
}

// handlerAlias applies the given alias mapping to the rule name and labels.
// If no alias is provided for a name or label, the original name is used.
func (br *baseRule) handlerAlias(alias map[string]string) {
	br.alias = alias

	if v, ok := alias[br.name]; ok {
		br.nameAlias = v
	} else {
		br.nameAlias = br.name
	}

	_labelsAlias := make(map[string]string, len(br.labels))
	_labelsAliasLength := make(map[string]int, len(br.labels))
	for label, length := range br.labels {
		if v, ok := alias[label]; ok {
			_labelsAlias[label] = v
			_labelsAliasLength[v] = length
		} else {
			_labelsAlias[label] = label
			_labelsAliasLength[label] = length
		}
	}

	br.labelsAlias = _labelsAlias
	br.labelsAliasLength = _labelsAliasLength
}
