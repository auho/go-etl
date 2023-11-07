package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type TagDataRulesTable struct {
	table
	tag *model.TagDataRules
}

func NewTagDataRulesTable(tag *model.TagDataRules) *TagDataRulesTable {
	t := &TagDataRulesTable{}
	t.tag = tag

	t.build()

	return t
}

func (t *TagDataRulesTable) build() {
	t.initCommand(t.tag.TableName())
	t.Command.AddPkInt(t.tag.GetIdName())

	NewDataTable(t.tag.GetData()).BuildForTag(t.Command)
	for _, rule := range t.tag.GetRules() {
		NewRuleTable(rule).BuildForTag(t.Command)
	}

	t.execCommand()
	t.execRowsCommand(t.tag)
}

func (t *TagDataRulesTable) WithCommand(fn func(command *tablestructure.Command)) *TagDataRulesTable {
	t.commandFun = fn

	return t
}
