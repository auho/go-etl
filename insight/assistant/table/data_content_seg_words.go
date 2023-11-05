package table

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
)

type DataContentSegWordsTable struct {
	table
	dataContentSegWords *model.DataContentSegWords
}

func NewDataContentSegWordsTable(d *model.DataContentSegWords) *DataContentSegWordsTable {
	t := &DataContentSegWordsTable{}
	t.dataContentSegWords = d

	t.buildSegWords()

	return t
}

func (t *DataContentSegWordsTable) buildSegWords() {
	t.initCommand(t.dataContentSegWords.TableName())

	t.command.AddPkInt("id")
	t.command.AddKeyBigInt(t.dataContentSegWords.GetData().GetIdName())
	t.command.AddStringWithLength(t.dataContentSegWords.GetContentName(), t.dataContentSegWords.GetContentLength())
	t.command.AddStringWithLength(t.dataContentSegWords.WordName(), 30)
	t.command.AddStringWithLength(t.dataContentSegWords.FlagName(), 5)
	t.command.AddInt(t.dataContentSegWords.NumName())
}
