package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
)

type extra struct {
	model assistant.Rawer
}

func (e *extra) DmlTable() *dml.Table {
	return dml.NewTable(e.model.TableName())
}

func (e *extra) Truncate() error {
	return e.model.GetDB().Truncate(e.model.TableName())
}

func (e *extra) CopyBuild(dst assistant.Rawer) error {
	return e.model.GetDB().DropAndCopy(e.model.TableName(), dst.TableName())
}

func (e *extra) InsertWithTable(table dml.Tabler) (string, error) {
	return table.Insert(e.model.TableName(), e.model.GetDB())
}

func (e *extra) InsertWithTableFiled(table dml.Tabler, fields []string) (string, error) {
	return table.InsertWithField(e.model.TableName(), fields, e.model.GetDB())
}
