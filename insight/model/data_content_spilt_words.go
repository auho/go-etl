package model

import (
	"fmt"
)

type DataContentSpiltWords struct {
	data          *Data
	contentName   string
	contentLength int
}

func NewDataContentSpiltWords(data *Data, contentName string, contentLength int) *DataContentSpiltWords {
	dcsw := &DataContentSpiltWords{}
	dcsw.data = data
	dcsw.contentName = contentName
	dcsw.contentLength = contentLength

	return dcsw
}

func (dcsw *DataContentSpiltWords) GetData() *Data {
	return dcsw.data
}

func (dcsw *DataContentSpiltWords) GetContentName() string {
	return dcsw.contentName
}

func (dcsw *DataContentSpiltWords) GetContentLength() int {
	return dcsw.contentLength
}

func (dcsw *DataContentSpiltWords) TableName() string {
	return fmt.Sprintf("%s_%s_%s_%s", NameTag, dcsw.data.GetName(), dcsw.contentName, NameSpiltWords)
}

func (dcsw *DataContentSpiltWords) WordName() string {
	return NameWord
}
