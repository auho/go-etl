package buildtable

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/ddl/command/mysql"
)

type RuleTable struct {
	table
	rule assistant.Rule
}

func NewRuleTable(rule assistant.Rule, opts ...TableOption) *RuleTable {
	t := &RuleTable{}
	t.rule = rule
	t.db = rule.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *RuleTable) build() {
	t.initCommand(t.rule.TableName())
	t.Command.AddPKInt(t.rule.IDName())

	t.BuildLabels(t.Command)

	var keywordFiled *mysql.Field
	if t.rule.Config().AllowKeywordDuplicate() {
		keywordFiled = t.Command.AddStringWithLength(t.rule.KeywordName(), t.rule.KeywordLength())
	} else {
		keywordFiled = t.Command.AddUniqueString(t.rule.KeywordName(), t.rule.KeywordLength())
	}

	keywordFiled.SetCollateUtf8mb4Bin()
	t.Command.AddInt(t.rule.KeywordLenName())

	t.execRawCommandFunc(t.rule)

	t.Command.AddTimestamp("ctime", true, true)
}

// BuildLabels
// labels
func (t *RuleTable) BuildLabels(command *schema.Command) {
	command.AddStringWithLength(t.rule.GetName(), t.rule.NameLength())

	for label, length := range t.rule.Labels() {
		command.AddStringWithLength(label, length)
	}
}

// BuildLabelsForWhole
// labels for whole
func (t *RuleTable) BuildLabelsForWhole(command *schema.Command, length int) {
	command.AddStringWithLength(t.rule.GetName(), length)

	for label := range t.rule.Labels() {
		command.AddStringWithLength(label, length)
	}
}

// BuildTags
// tags
func (t *RuleTable) BuildTags(command *schema.Command) {
	t.BuildLabels(command)
	command.AddStringWithLength(t.rule.KeywordName(), t.rule.KeywordLength())
}

// BuildTagsForWhole
// tags for whole
func (t *RuleTable) BuildTagsForWhole(command *schema.Command, length int) {
	t.BuildLabelsForWhole(command, length)
	command.AddStringWithLength(t.rule.KeywordName(), length)
}

// BuildForTag
// for tag
func (t *RuleTable) BuildForTag(command *schema.Command) {
	t.BuildTags(command)

	command.AddInt(t.rule.KeywordNumName())
	command.AddInt(t.rule.KeywordAmountName())
	command.AddInt(t.rule.LabelNumName())
}

func (t *RuleTable) WithCommand(fn func(*schema.Command)) *RuleTable {
	fn(t.Command)

	return t
}
