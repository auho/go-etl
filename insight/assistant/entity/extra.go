package entity

import (
	"context"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/schema/alter"
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

// extra provides common database operations as an embedded mixin.
// It wraps an assistant.Raw entity to offer DDL/DML helpers such as
// table creation, insertion, deletion, truncation, and schema copying.
type extra struct {
	raw assistant.Raw // the underlying entity providing table name and DB connection
}

// DMLTable returns a DML table reference for the underlying entity's table.
func (e *extra) DMLTable() *dml.Table {
	return dml.NewTable(e.raw.TableName())
}

// AlterTable builds and executes ALTER TABLE SQL commands via the model table builder.
// The optional fn parameter allows customizing the schema command before execution.
func (e *extra) AlterTable(fn func(*schema.Command)) ([]string, error) {
	at := alter.NewModelTable(e.raw).WithCommand(fn)
	return at.BuildAndReturnSQL()
}

// InsertWholeWithTable inserts all rows from the given table into the underlying entity's table.
func (e *extra) InsertWholeWithTable(table dml.Tabler) (string, error) {
	return table.Insert(e.raw.TableName(), e.raw.DB())
}

// InsertWithTable inserts selected fields from the given table into the underlying entity's table.
func (e *extra) InsertWithTable(table dml.Tabler) (string, error) {
	return table.InsertWithFields(e.raw.TableName(), table.GetSelectFields(), e.raw.DB())
}

// InsertWithTableField inserts the specified fields from the given table into the underlying entity's table.
func (e *extra) InsertWithTableField(table dml.Tabler, fields []string) (string, error) {
	return table.InsertWithFields(e.raw.TableName(), fields, e.raw.DB())
}

// GetTableColumns returns the column names of the underlying entity's table.
func (e *extra) GetTableColumns() ([]string, error) {
	return e.raw.DB().GetTableColumns(context.TODO(), e.raw.TableName())
}

// RecreateFromSql drops the table and recreates it using the provided SQL statement.
func (e *extra) RecreateFromSql(sql string) error {
	err := e.Drop()
	if err != nil {
		return err
	}

	return e.ExecSql(sql)
}

// Truncate truncates the underlying entity's table.
func (e *extra) Truncate() error {
	return e.raw.DB().Truncate(context.TODO(), e.raw.TableName())
}

// Drop drops the underlying entity's table.
func (e *extra) Drop() error {
	return e.raw.DB().Drop(context.TODO(), e.raw.TableName())
}

// CopyStructure copies the table structure (schema only, no data) from the underlying entity to the destination.
func (e *extra) CopyStructure(dst assistant.Raw) error {
	err := e.raw.DB().Drop(context.TODO(), dst.TableName())
	if err != nil {
		return err
	}

	return e.raw.DB().CopyStructure(context.TODO(), e.raw.TableName(), dst.TableName())
}

// CopyStructureAndData copies both the table structure and all data from the underlying entity to the destination.
func (e *extra) CopyStructureAndData(dst assistant.Raw) error {
	err := e.CopyStructure(dst)
	if err != nil {
		return err
	}

	return e.raw.DB().CopyData(context.TODO(), e.raw.TableName(), dst.TableName())
}

// RawSql executes a raw SQL query and scans the result into dst.
func (e *extra) RawSql(dst any, sql string, v ...any) error {
	return e.raw.DB().GormDB().WithContext(context.TODO()).Raw(sql, v...).Scan(dst).Error
}

// ExecSql executes a raw SQL statement (e.g., INSERT, UPDATE, DELETE).
func (e *extra) ExecSql(sql string, v ...any) error {
	return e.raw.DB().GormDB().WithContext(context.TODO()).Exec(sql, v...).Error
}
