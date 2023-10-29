package model

import (
	"fmt"
)

type DataContentSegWords struct {
	data          *Data
	contentName   string
	contentLength int
}

func NewDataContentSegWords(data *Data, contentName string, contentLength int) *DataContentSegWords {
	dcsw := &DataContentSegWords{}
	dcsw.data = data
	dcsw.contentName = contentName
	dcsw.contentLength = contentLength

	return dcsw
}

func (dcsw *DataContentSegWords) GetData() *Data {
	return dcsw.data
}

func (dcsw *DataContentSegWords) GetContentName() string {
	return dcsw.contentName
}

func (dcsw *DataContentSegWords) GetContentLength() int {
	return dcsw.contentLength
}

func (dcsw *DataContentSegWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dcsw.data.GetName(), dcsw.contentName, NameSegWords)
}

func (dcsw *DataContentSegWords) WordName() string {
	return NameWord
}

func (dcsw *DataContentSegWords) FlagName() string {
	return NameFlag
}

func (dcsw *DataContentSegWords) NumName() string {
	return NameNum
}
