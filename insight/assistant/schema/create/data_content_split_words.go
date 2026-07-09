package create

import (
	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

type DataContentSplitWordsTable struct {
	table
	dataContentSplitWords *entity.DataContentSplitWords
}

func NewDataContentSplitWordsTable(d *entity.DataContentSplitWords, opts ...TableOption) *DataContentSplitWordsTable {
	t := &DataContentSplitWordsTable{}
	t.dataContentSplitWords = d
	t.db = d.DB()

	t.options(opts)
	t.build()

	return t
}

func (t *DataContentSplitWordsTable) build() {
	t.initCommand(t.dataContentSplitWords.TableName())

	t.Command.AddPKInt(t.dataContentSplitWords.IDName())
	t.Command.AddKeyBigInt(t.dataContentSplitWords.GetData().IDName())
	t.Command.AddStringWithLength(t.dataContentSplitWords.WordName(), 30)

	t.execRawCommandFunc(t.dataContentSplitWords)
}

func (t *DataContentSplitWordsTable) WithCommand(fn func(*schema.Command)) *DataContentSplitWordsTable {
	fn(t.Command)

	return t
}
