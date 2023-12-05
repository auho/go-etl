package model

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/altertable"
)

type extra struct {
	model assistant.Rawer
}

func (e *extra) DmlTable() *dml.Table {
	return dml.NewTable(e.model.TableName())
}

func (e *extra) AlterTable(fn func(*tablestructure.Command)) ([]string, error) {
	at := altertable.NewModelTable(e.model).WithCommand(fn)
	return at.BuildAffixSql()
}

func (e *extra) InsertFromTable(table dml.Tabler) (string, error) {
	return table.Insert(e.model.TableName(), e.model.GetDB())
}

func (e *extra) InsertWithTable(table dml.Tabler) (string, error) {
	return table.InsertWithField(e.model.TableName(), table.GetSelectFields(), e.model.GetDB())
}

func (e *extra) InsertWithTableFiled(table dml.Tabler, fields []string) (string, error) {
	return table.InsertWithField(e.model.TableName(), fields, e.model.GetDB())
}

func (e *extra) GetTableColumns() ([]string, error) {
	return e.model.GetDB().GetTableColumns(e.model.TableName())
}

func (e *extra) Truncate() error {
	return e.model.GetDB().Truncate(e.model.TableName())
}

func (e *extra) CopyBuild(dst assistant.Rawer) error {
	return e.model.GetDB().DropAndCopy(e.model.TableName(), dst.TableName())
}
