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
	base assistant.Raw
}

func (e *extra) DMLTable() *dml.Table {
	return dml.NewTable(e.base.TableName())
}

func (e *extra) AlterTable(fn func(*schema.Command)) ([]string, error) {
	at := alter.NewModelTable(e.base).WithCommand(fn)
	return at.BuildAndReturnSQL()
}

func (e *extra) InsertWholeWithTable(table dml.Tabler) (string, error) {
	return table.Insert(e.base.TableName(), e.base.DB())
}

func (e *extra) InsertWithTable(table dml.Tabler) (string, error) {
	return table.InsertWithFields(e.base.TableName(), table.GetSelectFields(), e.base.DB())
}

func (e *extra) InsertWithTableField(table dml.Tabler, fields []string) (string, error) {
	return table.InsertWithFields(e.base.TableName(), fields, e.base.DB())
}

func (e *extra) GetTableColumns() ([]string, error) {
	return e.base.DB().GetTableColumns(context.TODO(), e.base.TableName())
}

func (e *extra) RecreateFromSql(sql string) error {
	err := e.Drop()
	if err != nil {
		return err
	}

	return e.ExecSql(sql)
}

func (e *extra) Drop() error {
	return e.base.DB().Drop(context.TODO(), e.base.TableName())
}

func (e *extra) Truncate() error {
	return e.base.DB().Truncate(context.TODO(), e.base.TableName())
}

func (e *extra) CopyBuild(dst assistant.Raw) error {
	err := dst.DB().Drop(context.TODO(), dst.TableName())
	if err != nil {
		return err
	}

	return e.base.DB().CopyStructure(context.TODO(), e.base.TableName(), dst.TableName())
}

func (e *extra) CopyBuildAndData(dst assistant.Raw) error {
	err := e.CopyBuild(dst)
	if err != nil {
		return err
	}

	return e.base.DB().GormDB().WithContext(context.TODO()).Exec(
		fmt.Sprintf("INSERT INTO %s SELECT * FROM %s", dst.TableName(), e.base.TableName()),
	).Error
}

func (e *extra) RawSqlAndScan(dst any, sql string, v ...any) error {
	return e.base.DB().GormDB().WithContext(context.TODO()).Raw(sql, v...).Scan(dst).Error
}

func (e *extra) ExecSql(sql string, v ...any) error {
	return e.base.DB().GormDB().WithContext(context.TODO()).Exec(sql, v...).Error
}
