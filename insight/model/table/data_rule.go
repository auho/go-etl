package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type DataRuleTable struct {
	table
	dataRule *model.DataRule
}

func NewDataRuleTable(dataRule *model.DataRule) *DataRuleTable {
	rt := &DataRuleTable{}
	rt.dataRule = dataRule

	rt.buildDataRule()

	return rt
}

func (rt *DataRuleTable) buildDataRule() {
	rt.initTable(rt.dataRule.TableName())
	rt.AddPkInt("id")

	rt.AddStringWithLength(rt.dataRule.GetRule().GetName(), rt.dataRule.GetRule().GetLength())

	for label, length := range rt.dataRule.GetRule().GetLabels() {
		rt.table.AddStringWithLength(label, length)
	}

	rt.AddUniqueString(rt.dataRule.GetRule().Keyword(), rt.dataRule.GetRule().GetKeywordLength())
	rt.AddInt(rt.dataRule.GetRule().KeywordLen())
}
