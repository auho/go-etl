package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
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

	t.buildRuleLabels(t.Command)
	keywordFiled := t.Command.AddUniqueString(t.rule.KeywordName(), t.rule.GetNameLength())
	keywordFiled.SetCollateUtf8mb4Bin()
	t.Command.AddInt(t.rule.KeywordLenName())

	t.execCommand()
	t.execRowsCommand(t.rule)

	t.Command.AddTimestamp("ctime", true, true)
}

func (t *RuleTable) buildRuleLabels(command *tablestructure.Command) {
	command.AddStringWithLength(t.rule.GetName(), t.rule.GetNameLength())

	for label, length := range t.rule.GetLabels() {
		command.AddStringWithLength(label, length)
	}
}

func (t *RuleTable) BuildForTag(command *tablestructure.Command) {
	t.buildRuleLabels(command)
	command.AddStringWithLength(t.rule.KeywordName(), t.rule.GetNameLength())
	command.AddInt(t.rule.KeywordNumName())
}
