package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
)

type TagDataRuleTable struct {
	table
	tag *model.TagDataRule
}

func NewTagDataRuleTable(tag *model.TagDataRule, opts ...TableOption) *TagDataRuleTable {
	t := &TagDataRuleTable{}
	t.tag = tag
	t.db = tag.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *TagDataRuleTable) build() {
	t.initCommand(t.tag.TableName())
	t.Command.AddPkInt(t.tag.GetIdName())

	NewDataTable(t.tag.GetData()).BuildForTag(t.Command)
	NewRuleTable(t.tag.GetRule()).BuildForTag(t.Command)

	t.execCommand()
	t.execRawCommand(t.tag)
}
