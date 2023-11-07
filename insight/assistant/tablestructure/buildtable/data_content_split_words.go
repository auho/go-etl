package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type DataContentSpiltWordsTable struct {
	table
	dataContentSpiltWords *model.DataContentSpiltWords
}

func NewDataContentSpiltWordsTable(d *model.DataContentSpiltWords) *DataContentSpiltWordsTable {
	t := &DataContentSpiltWordsTable{}
	t.dataContentSpiltWords = d

	t.build()

	return t
}

func (t *DataContentSpiltWordsTable) build() {
	t.initCommand(t.dataContentSpiltWords.TableName())

	t.Command.AddPkInt(t.dataContentSpiltWords.GetIdName())
	t.Command.AddKeyBigInt(t.dataContentSpiltWords.GetData().GetIdName())
	t.Command.AddStringWithLength(t.dataContentSpiltWords.GetContentName(), t.dataContentSpiltWords.GetContentLength())
	t.Command.AddStringWithLength(t.dataContentSpiltWords.WordName(), 30)

	t.execCommand()
	t.execRowsCommand(t.dataContentSpiltWords)
}

func (t *DataContentSpiltWordsTable) WithCommand(fn func(command *tablestructure.Command)) *DataContentSpiltWordsTable {
	t.commandFun = fn

	return t
}
