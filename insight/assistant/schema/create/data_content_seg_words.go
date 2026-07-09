package create

import (
	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

type DataContentSegWordsTable struct {
	table
	dataContentSegWords *entity.DataContentSegWords
}

func NewDataContentSegWordsTable(d *entity.DataContentSegWords, opts ...TableOption) *DataContentSegWordsTable {
	t := &DataContentSegWordsTable{}
	t.dataContentSegWords = d
	t.db = d.DB()

	t.options(opts)
	t.build()

	return t
}

func (t *DataContentSegWordsTable) build() {
	t.initCommand(t.dataContentSegWords.TableName())

	t.Command.AddPKInt(t.dataContentSegWords.IDName())
	t.Command.AddKeyBigInt(t.dataContentSegWords.GetData().IDName())
	t.Command.AddStringWithLength(t.dataContentSegWords.WordName(), 30)
	t.Command.AddStringWithLength(t.dataContentSegWords.FlagName(), 5)
	t.Command.AddInt(t.dataContentSegWords.NumName())

	t.execRawCommandFunc(t.dataContentSegWords)
}

func (t *DataContentSegWordsTable) WithCommand(fn func(*schema.Command)) *DataContentSegWordsTable {
	fn(t.Command)

	return t
}
