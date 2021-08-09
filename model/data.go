package model

import (
	"fmt"

	goEtl "github.com/auho/go-etl"
)

type Data struct {
	name string
}

func NewData(name string) *Data {
	d := &Data{}
	d.name = name

	return d
}

func (d *Data) getName() string {
	return d.name
}

func (d *Data) tableName() string {
	return d.name
}

func (d *Data) tagTableName(tagName string) string {
	return fmt.Sprintf("%s_%s_%s", goEtl.TagTableNamePrefix, d.name, tagName)
}
