package model

import (
	"fmt"

	"github.com/auho/go-etl/tool"
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
	return d.name
}

func (d *Data) TagTableName(tagName string) string {
	return fmt.Sprintf("%s_%s_%s", tool.TagTableNamePrefix, d.name, tagName)
}

func (d *Data) RuleTableName(n string) string {
	return fmt.Sprintf("%s_%s_%s", tool.RuleTableNamePrefix, d.name, n)
}
