package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/schema"
)

type DataContentSplitWordsTable struct {
	table
	dataContentSplitWords *model.DataContentSplitWords
}

func NewDataContentSplitWordsTable(d *model.DataContentSplitWords, opts ...TableOption) *DataContentSplitWordsTable {
	t := &DataContentSplitWordsTable{}
	t.dataContentSplitWords = d
	t.db = d.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *DataContentSplitWordsTable) build() {
	t.initCommand(t.dataContentSplitWords.TableName())

	t.Command.AddPKInt(t.dataContentSplitWords.GetIDName())
	t.Command.AddKeyBigInt(t.dataContentSplitWords.GetData().GetIDName())
	t.Command.AddStringWithLength(t.dataContentSplitWords.WordName(), 30)

	t.execRawCommandFunc(t.dataContentSplitWords)
}

func (t *DataContentSplitWordsTable) WithCommand(fn func(*schema.Command)) *DataContentSplitWordsTable {
	fn(t.Command)

	return t
}
