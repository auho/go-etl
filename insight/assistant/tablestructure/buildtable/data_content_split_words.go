package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
)

type DataContentSpiltWordsTable struct {
	table
	dataContentSpiltWords *model.DataContentSpiltWords
}

func NewDataContentSpiltWordsTable(d *model.DataContentSpiltWords, opts ...TableOption) *DataContentSpiltWordsTable {
	t := &DataContentSpiltWordsTable{}
	t.dataContentSpiltWords = d
	t.db = d.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *DataContentSpiltWordsTable) build() {
	t.initCommand(t.dataContentSpiltWords.TableName())

	t.Command.AddPkInt(t.dataContentSpiltWords.GetIdName())
	t.Command.AddKeyBigInt(t.dataContentSpiltWords.GetData().GetIdName())
	t.Command.AddStringWithLength(t.dataContentSpiltWords.WordName(), 30)

	t.execCommand()
	t.execRowsCommand(t.dataContentSpiltWords)
}
