package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type RuleTable struct {
	table
	rule model.Ruler
}

func NewRuleTable(rule model.Ruler) *RuleTable {
	rt := &RuleTable{}
	rt.rule = rule

	rt.buildRule()

	return rt
}

func (rt *RuleTable) buildRule() {
	rt.initTable(rt.rule.TableName())
	rt.AddPkInt(rt.rule.GetIdName())

	rt.AddStringWithLength(rt.rule.GetName(), rt.rule.GetNameLength())

	for label, length := range rt.rule.GetLabels() {
		rt.table.AddStringWithLength(label, length)
	}

	rt.AddUniqueString(rt.rule.KeywordName(), rt.rule.GetNameLength())
	rt.AddInt(rt.rule.KeywordLenName())
	rt.AddTimestamp("ctime", true, true)

}
