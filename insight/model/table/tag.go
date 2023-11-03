package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type TagTable struct {
	table
	tag *model.Tag
}

func NewTagTable(tag *model.Tag) *TagTable {
	tt := &TagTable{}
	tt.tag = tag

	tt.buildTag()

	return tt
}

func (tt *TagTable) buildTag() {
	tt.initTable(tt.tag.TableName())
	tt.AddPkInt("id")

	tt.AddKeyBigInt(tt.tag.GetData().GetIdName())
	tt.AddStringWithLength(tt.tag.GetRule().GetName(), tt.tag.GetRule().GetNameLength())

	for label, length := range tt.tag.GetRule().GetLabels() {
		tt.table.AddStringWithLength(label, length)
	}

	tt.AddStringWithLength(tt.tag.GetRule().KeywordName(), tt.tag.GetRule().GetKeywordLength())
	tt.AddInt(tt.tag.KeywordNumName())
}
