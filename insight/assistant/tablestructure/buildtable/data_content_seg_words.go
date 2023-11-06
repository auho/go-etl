package buildtable

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

	t.build()

	return t
}

func (t *DataContentSegWordsTable) build() {
	t.initCommand(t.dataContentSegWords.TableName())

	t.Command.AddPkInt("id")
	t.Command.AddKeyBigInt(t.dataContentSegWords.GetData().GetIdName())
	t.Command.AddStringWithLength(t.dataContentSegWords.GetContentName(), t.dataContentSegWords.GetContentLength())
	t.Command.AddStringWithLength(t.dataContentSegWords.WordName(), 30)
	t.Command.AddStringWithLength(t.dataContentSegWords.FlagName(), 5)
	t.Command.AddInt(t.dataContentSegWords.NumName())
}
