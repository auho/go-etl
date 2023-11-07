package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
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

	t.Command.AddPkInt(t.dataContentSegWords.GetIdName())
	t.Command.AddKeyBigInt(t.dataContentSegWords.GetData().GetIdName())
	t.Command.AddStringWithLength(t.dataContentSegWords.GetContentName(), t.dataContentSegWords.GetContentLength())
	t.Command.AddStringWithLength(t.dataContentSegWords.WordName(), 30)
	t.Command.AddStringWithLength(t.dataContentSegWords.FlagName(), 5)
	t.Command.AddInt(t.dataContentSegWords.NumName())

	t.execCommand()
	t.execRowsCommand(t.dataContentSegWords)
}

func (t *DataContentSegWordsTable) WithCommand(fn func(command *tablestructure.Command)) *DataContentSegWordsTable {
	t.commandFun = fn

	return t
}
