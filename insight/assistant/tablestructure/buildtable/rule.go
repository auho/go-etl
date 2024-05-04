package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/ddl/command/mysql"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type RuleTable struct {
	table
	rule assistant.Ruler
}

func NewRuleTable(rule assistant.Ruler, opts ...TableOption) *RuleTable {
	t := &RuleTable{}
	t.rule = rule
	t.db = rule.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *RuleTable) build() {
	t.initCommand(t.rule.TableName())
	t.Command.AddPkInt(t.rule.GetIdName())

	t.BuildLabels(t.Command)

	var keywordFiled *mysql.Field
	if t.rule.Config().AllowKeywordDuplicate() {
		keywordFiled = t.Command.AddStringWithLength(t.rule.KeywordName(), t.rule.GetKeywordLength())
	} else {
		keywordFiled = t.Command.AddUniqueString(t.rule.KeywordName(), t.rule.GetKeywordLength())
	}

	keywordFiled.SetCollateUtf8mb4Bin()
	t.Command.AddInt(t.rule.KeywordLenName())

	t.execRawCommandFunc(t.rule)

	t.Command.AddTimestamp("ctime", true, true)
}

// BuildLabels
// labels
func (t *RuleTable) BuildLabels(command *tablestructure.Command) {
	command.AddStringWithLength(t.rule.GetName(), t.rule.GetNameLength())

	for label, length := range t.rule.GetLabels() {
		command.AddStringWithLength(label, length)
	}
}

// BuildLabelsForWhole
// labels for whole
func (t *RuleTable) BuildLabelsForWhole(command *tablestructure.Command, length int) {
	command.AddStringWithLength(t.rule.GetName(), length)

	for label := range t.rule.GetLabels() {
		command.AddStringWithLength(label, length)
	}
}

// BuildTags
// tags
func (t *RuleTable) BuildTags(command *tablestructure.Command) {
	t.BuildLabels(command)
	command.AddStringWithLength(t.rule.KeywordName(), t.rule.GetKeywordLength())
}

// BuildTagsForWhole
// tags for whole
func (t *RuleTable) BuildTagsForWhole(command *tablestructure.Command, length int) {
	t.BuildLabelsForWhole(command, length)
	command.AddStringWithLength(t.rule.KeywordName(), length)
}

// BuildForTag
// for tag
func (t *RuleTable) BuildForTag(command *tablestructure.Command) {
	t.BuildTags(command)

	command.AddInt(t.rule.KeywordNumName())
	command.AddInt(t.rule.KeywordAmountName())
	command.AddInt(t.rule.LabelNumName())
}

func (t *RuleTable) WithCommand(fn func(*tablestructure.Command)) *RuleTable {
	fn(t.Command)

	return t
}
