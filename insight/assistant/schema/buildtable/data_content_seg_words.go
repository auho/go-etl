package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/schema"
)

type DataContentSegWordsTable struct {
	table
	dataContentSegWords *model.DataContentSegWords
}

func NewDataContentSegWordsTable(d *model.DataContentSegWords, opts ...TableOption) *DataContentSegWordsTable {
	t := &DataContentSegWordsTable{}
	t.dataContentSegWords = d
	t.db = d.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *DataContentSegWordsTable) build() {
	t.initCommand(t.dataContentSegWords.TableName())

	t.Command.AddPkInt(t.dataContentSegWords.GetIDName())
	t.Command.AddKeyBigInt(t.dataContentSegWords.GetData().GetIDName())
	t.Command.AddStringWithLength(t.dataContentSegWords.WordName(), 30)
	t.Command.AddStringWithLength(t.dataContentSegWords.FlagName(), 5)
	t.Command.AddInt(t.dataContentSegWords.NumName())

	t.execRawCommandFunc(t.dataContentSegWords)
}

func (t *DataContentSegWordsTable) WithCommand(fn func(*schema.Command)) *DataContentSegWordsTable {
	fn(t.Command)

	return t
}
