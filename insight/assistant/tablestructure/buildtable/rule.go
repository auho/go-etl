package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type RuleTable struct {
	table
	rule assistant.Ruler
}

func NewRuleTable(rule assistant.Ruler) *RuleTable {
	t := &RuleTable{}
	t.rule = rule

	t.build()

	return t
}

func (t *RuleTable) build() {
	t.initCommand(t.rule.TableName())
	t.Command.AddPkInt(t.rule.GetIdName())

	t.buildRule(t.Command)

	t.Command.AddInt(t.rule.KeywordLenName())
	t.execRowsCommand(t.rule)
	t.Command.AddTimestamp("ctime", true, true)
}

func (t *RuleTable) buildRule(command *tablestructure.Command) {
	command.AddStringWithLength(t.rule.GetName(), t.rule.GetNameLength())

	for label, length := range t.rule.GetLabels() {
		command.AddStringWithLength(label, length)
	}

	keywordFiled := t.Command.AddUniqueString(t.rule.KeywordName(), t.rule.GetNameLength())
	keywordFiled.SetCollateUtf8mb4Bin()
}

func (t *RuleTable) BuildForTag(command *tablestructure.Command) {
	t.buildRule(command)
	command.AddInt(t.rule.KeywordNumName())
}
