package table

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
)

type TagDataRuleTable struct {
	table
	tag *model.TagDataRule
}

func NewTagDataRuleTable(tag *model.TagDataRule) *TagDataRuleTable {
	t := &TagDataRuleTable{}
	t.tag = tag

	t.buildTag()

	return t
}

func (t *TagDataRuleTable) buildTag() {
	t.initCommand(t.tag.TableName())
	t.command.AddPkInt(t.tag.GetIdName())

	NewDataTable(t.tag.GetData()).BuildDataForTag(t.command)
	NewRuleTable(t.tag.GetRule()).BuildRuleForTag(t.command)
}
