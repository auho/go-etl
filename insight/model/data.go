package model

import (
	"fmt"
)

type Data struct {
	name   string
	idName string
}

func NewData(name string, idName string) *Data {
	d := &Data{}
	d.name = name
	d.idName = idName

	return d
}

func (d *Data) GetName() string {
	return d.name
}

func (d *Data) GetIdName() string {
	return d.idName
}

func (d *Data) TableName() string {
	return fmt.Sprintf("%s_%s", NameData, d.name)
}
