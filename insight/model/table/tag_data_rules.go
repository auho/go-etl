package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type TagDataRulesTable struct {
	table
	tag *model.TagDataRules
}

func NewTagDataRulesTable(tag *model.TagDataRules) *TagDataRulesTable {
	t := &TagDataRulesTable{}
	t.tag = tag

	t.buildTag()

	return t
}

func (t *TagDataRulesTable) buildTag() {
	t.initCommand(t.tag.TableName())
	t.command.AddPkInt(t.tag.GetIdName())

	NewDataTable(t.tag.GetData()).BuildDataForTag(t.command)
	for _, rule := range t.tag.GetRules() {
		NewRuleTable(rule).BuildRuleForTag(t.command)
	}
}
