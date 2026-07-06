package entity

import (
	"context"
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/schema/alter"
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

type extra struct {
	model assistant.Raw
}

func (e *extra) DMLTable() *dml.Table {
	return dml.NewTable(e.model.TableName())
}

func (e *extra) AlterTable(fn func(*schema.Command)) ([]string, error) {
	at := alter.NewModelTable(e.model).WithCommand(fn)
	return at.BuildAndReturnSQL()
}

func (e *extra) InsertWholeWithTable(table dml.Tabler) (string, error) {
	return table.Insert(e.model.TableName(), e.model.GetDB())
}

func (e *extra) InsertWithTable(table dml.Tabler) (string, error) {
	return table.InsertWithFields(e.model.TableName(), table.GetSelectFields(), e.model.GetDB())
}

func (e *extra) InsertWithTableField(table dml.Tabler, fields []string) (string, error) {
	return table.InsertWithFields(e.model.TableName(), fields, e.model.GetDB())
}

func (e *extra) GetTableColumns() ([]string, error) {
	return e.model.GetDB().GetTableColumns(context.TODO(), e.model.TableName())
}

func (e *extra) RecreateFromSql(sql string) error {
	err := e.Drop()
	if err != nil {
		return err
	}

	return e.ExecSql(sql)
}

func (e *extra) Drop() error {
	return e.model.GetDB().Drop(context.TODO(), e.model.TableName())
}

func (e *extra) Truncate() error {
	return e.model.GetDB().Truncate(context.TODO(), e.model.TableName())
}

func (e *extra) CopyBuild(dst assistant.Raw) error {
	err := dst.GetDB().Drop(context.TODO(), dst.TableName())
	if err != nil {
		return err
	}

	return e.model.GetDB().CopyStructure(context.TODO(), e.model.TableName(), dst.TableName())
}

func (e *extra) CopyBuildAndData(dst assistant.Raw) error {
	err := e.CopyBuild(dst)
	if err != nil {
		return err
	}

	return e.model.GetDB().GormDB().WithContext(context.TODO()).Exec(
		fmt.Sprintf("INSERT INTO %s SELECT * FROM %s", dst.TableName(), e.model.TableName()),
	).Error
}

func (e *extra) RawSqlAndScan(dst any, sql string, v ...any) error {
	return e.model.GetDB().GormDB().WithContext(context.TODO()).Raw(sql, v...).Scan(dst).Error
}

func (e *extra) ExecSql(sql string, v ...any) error {
	return e.model.GetDB().GormDB().WithContext(context.TODO()).Exec(sql, v...).Error
}
