package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type TagDataRuleTable struct {
	table
	tag *model.TagDataRule
}

func NewTagDataRuleTable(tag *model.TagDataRule) *TagDataRuleTable {
	t := &TagDataRuleTable{}
	t.tag = tag

	t.build()

	return t
}

func (t *TagDataRuleTable) build() {
	t.initCommand(t.tag.TableName())
	t.Command.AddPkInt(t.tag.GetIdName())

	NewDataTable(t.tag.GetData()).BuildForTag(t.Command)
	NewRuleTable(t.tag.GetRule()).BuildForTag(t.Command)

	t.execCommand()
	t.execRowsCommand(t.tag)
}

func (t *TagDataRuleTable) WithCommand(fn func(command *tablestructure.Command)) *TagDataRuleTable {
	t.commandFun = fn

	return t
}
