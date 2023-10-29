package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type DataContentSpiltWordsTable struct {
	table
	dataContentSpiltWords *model.DataContentSpiltWords
}

func NewDataContentSpiltWordsTable(d *model.DataContentSpiltWords) *DataContentSpiltWordsTable {
	swt := &DataContentSpiltWordsTable{}
	swt.dataContentSpiltWords = d

	swt.buildSpiltWords()

	return swt
}

func (dcswt *DataContentSpiltWordsTable) buildSpiltWords() {
	dcswt.initTable(dcswt.dataContentSpiltWords.TableName())

	dcswt.AddPkInt("id")
	dcswt.AddKeyBigInt(dcswt.dataContentSpiltWords.GetData().GetIdName())
	dcswt.AddStringWithLength(dcswt.dataContentSpiltWords.GetContentName(), dcswt.dataContentSpiltWords.GetContentLength())
	dcswt.AddStringWithLength(dcswt.dataContentSpiltWords.WordName(), 30)
}
