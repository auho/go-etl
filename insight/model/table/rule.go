package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type RuleTable struct {
	table
	rule model.Ruler
}

func NewRuleTable(rule model.Ruler) *RuleTable {
	t := &RuleTable{}
	t.rule = rule

	t.buildRule()

	return t
}

func (t *RuleTable) buildRule() {
	t.initCommand(t.rule.TableName())
	t.command.AddPkInt(t.rule.GetIdName())

	t.command.AddStringWithLength(t.rule.GetName(), t.rule.GetNameLength())

	for label, length := range t.rule.GetLabels() {
		t.command.AddStringWithLength(label, length)
	}

	t.command.AddUniqueString(t.rule.KeywordName(), t.rule.GetNameLength())
	t.command.AddInt(t.rule.KeywordLenName())
	t.command.AddTimestamp("ctime", true, true)

}

func (t *RuleTable) BuildRuleForTag(command *command) {
	command.AddStringWithLength(t.rule.GetName(), t.rule.GetNameLength())

	for label, length := range t.rule.GetLabels() {
		command.AddStringWithLength(label, length)
	}

	command.AddStringWithLength(t.rule.KeywordName(), t.rule.GetKeywordLength())
	command.AddInt(t.rule.KeywordNumName())
}
