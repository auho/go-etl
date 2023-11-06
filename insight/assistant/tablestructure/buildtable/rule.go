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

	t.buildRule()

	return t
}

func (t *RuleTable) buildRule() {
	t.initCommand(t.rule.TableName())
	t.Command.AddPkInt(t.rule.GetIdName())

	t.Command.AddStringWithLength(t.rule.GetName(), t.rule.GetNameLength())

	for label, length := range t.rule.GetLabels() {
		t.Command.AddStringWithLength(label, length)
	}

	t.Command.AddUniqueString(t.rule.KeywordName(), t.rule.GetNameLength())
	t.Command.AddInt(t.rule.KeywordLenName())
	t.Command.AddTimestamp("ctime", true, true)

}

func (t *RuleTable) BuildRuleForTag(command *tablestructure.Command) {
	command.AddStringWithLength(t.rule.GetName(), t.rule.GetNameLength())

	for label, length := range t.rule.GetLabels() {
		command.AddStringWithLength(label, length)
	}

	command.AddStringWithLength(t.rule.KeywordName(), t.rule.GetKeywordLength())
	command.AddInt(t.rule.KeywordNumName())
}
