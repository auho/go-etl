package create

import (
	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

type TagDataRuleTable struct {
	table
	tag *entity.TagDataRule
}

func NewTagDataRuleTable(tag *entity.TagDataRule, opts ...TableOption) *TagDataRuleTable {
	t := &TagDataRuleTable{}
	t.tag = tag
	t.db = tag.DB()

	t.options(opts)
	t.build()

	return t
}

func (t *TagDataRuleTable) build() {
	t.initCommand(t.tag.TableName())
	t.Command.AddPKInt(t.tag.IDName())

	NewDataTable(t.tag.GetData()).BuildForTag(t.Command)
	NewRuleTable(t.tag.GetRule()).BuildForTag(t.Command)

	t.execRawCommandFunc(t.tag)
}

func (t *TagDataRuleTable) WithCommand(fn func(*schema.Command)) *TagDataRuleTable {
	fn(t.Command)

	return t
}
