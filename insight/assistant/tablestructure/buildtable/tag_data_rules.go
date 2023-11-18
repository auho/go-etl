package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type TagDataRulesTable struct {
	table
	tag *model.TagDataRules
}

func NewTagDataRulesTable(tag *model.TagDataRules, opts ...TableOption) *TagDataRulesTable {
	t := &TagDataRulesTable{}
	t.tag = tag
	t.db = tag.GetDB()

	t.options(opts)
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

	t.execRawCommandFunc(t.tag)
}

func (t *TagDataRulesTable) WithCommand(fn func(*tablestructure.Command)) *TagDataRulesTable {
	fn(t.Command)

	return t
}
