package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
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

	t.execCommandFunc()
	t.execRawCommandFunc(t.tag)
}

func (t *TagDataRuleTable) WithCommand(fn func(*tablestructure.Command)) *TagDataRuleTable {
	fn(t.Command)

	return t
}
