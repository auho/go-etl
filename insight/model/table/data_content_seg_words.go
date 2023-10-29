package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type DataContentSegWordsTable struct {
	table
	dataContentSegWords *model.DataContentSegWords
}

func NewDataContentSegWordsTable(d *model.DataContentSegWords) *DataContentSegWordsTable {
	swt := &DataContentSegWordsTable{}
	swt.dataContentSegWords = d

	swt.buildSegWords()

	return swt
}

func (dcswt *DataContentSegWordsTable) buildSegWords() {
	dcswt.initTable(dcswt.dataContentSegWords.TableName())

	dcswt.AddPkInt("id")
	dcswt.AddKeyBigInt(dcswt.dataContentSegWords.GetData().GetIdName())
	dcswt.AddStringWithLength(dcswt.dataContentSegWords.GetContentName(), dcswt.dataContentSegWords.GetContentLength())
	dcswt.AddStringWithLength(dcswt.dataContentSegWords.WordName(), 30)
	dcswt.AddStringWithLength(dcswt.dataContentSegWords.FlagName(), 5)
	dcswt.AddInt(dcswt.dataContentSegWords.NumName())
}
