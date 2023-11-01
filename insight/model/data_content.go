package model

import (
	"fmt"
)

var _ Dataor = (*DataContent)(nil)

type DataContent struct {
	*Data
	contentName string
}

func NewDataContent(data *Data, contentName string) *DataContent {
	d := &DataContent{}
	d.Data = data
	d.contentName = contentName

	return d
}

func (d *DataContent) GetName() string {
	return fmt.Sprintf("%s_%s", d.name, d.contentName)
}
