package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type RuleTable struct {
	table
	rule *model.Rule
}

func NewRuleTable(rule *model.Rule) *RuleTable {
	rt := &RuleTable{}
	rt.rule = rule

	rt.buildRule()

	return rt
}

func (rt *RuleTable) buildRule() {
	rt.initTable(rt.rule.TableName())
	rt.AddPkInt("id")

	rt.AddStringWithLength(rt.rule.GetName(), rt.rule.GetLength())

	for label, length := range rt.rule.GetLabels() {
		rt.table.AddStringWithLength(label, length)
	}

	rt.AddUniqueString(rt.rule.Keyword(), rt.rule.GetLength())
	rt.AddInt(rt.rule.KeywordLen())
}
