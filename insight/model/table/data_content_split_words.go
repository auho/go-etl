package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type DataContentSpiltWordsTable struct {
	table
	dataContentSpiltWords *model.DataContentSpiltWords
}

func NewDataContentSpiltWordsTable(d *model.DataContentSpiltWords) *DataContentSpiltWordsTable {
	t := &DataContentSpiltWordsTable{}
	t.dataContentSpiltWords = d

	t.buildSpiltWords()

	return t
}

func (t *DataContentSpiltWordsTable) buildSpiltWords() {
	t.initCommand(t.dataContentSpiltWords.TableName())

	t.command.AddPkInt("id")
	t.command.AddKeyBigInt(t.dataContentSpiltWords.GetData().GetIdName())
	t.command.AddStringWithLength(t.dataContentSpiltWords.GetContentName(), t.dataContentSpiltWords.GetContentLength())
	t.command.AddStringWithLength(t.dataContentSpiltWords.WordName(), 30)
}
