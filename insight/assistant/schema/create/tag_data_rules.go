package create

import (
	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

type TagDataRulesTable struct {
	table
	tag *entity.TagDataRules
}

func NewTagDataRulesTable(tag *entity.TagDataRules, opts ...TableOption) *TagDataRulesTable {
	t := &TagDataRulesTable{}
	t.tag = tag
	t.db = tag.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *TagDataRulesTable) build() {
	t.initCommand(t.tag.TableName())
	t.Command.AddPKInt(t.tag.IDName())

	NewDataTable(t.tag.GetData()).BuildForTag(t.Command)
	for _, rule := range t.tag.GetRules() {
		NewRuleTable(rule).BuildForTag(t.Command)
	}

	t.execRawCommandFunc(t.tag)
}

func (t *TagDataRulesTable) WithCommand(fn func(*schema.Command)) *TagDataRulesTable {
	fn(t.Command)

	return t
}
