package entity

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/task/pipeline/database/etl"
)

var _ etl.CleanResource = (*CleanRows)(nil)

// CleanRows wraps the data cleaning lifecycle for a Rows entity.
// It provides access to the source (Raw entity), the rows being cleaned (Rows),
// and the deleted rows (Rows). This is used in the ETL pipeline's clean stage
// when working with raw source data rather than processed Data entities.
type CleanRows struct {
	base // embedded base for command hook and DB connection

	raw     assistant.Entity // source raw data entity
	rows    *Rows            // rows being cleaned
	deleted *Rows            // rows that were removed during cleaning
}

// NewCleanRows creates a new CleanRows for the given raw entity.
// A new Rows entity is created with the given name for the cleaned data output.
func NewCleanRows(newName string, raw assistant.Entity) *CleanRows {
	cd := &CleanRows{}
	cd.raw = raw
	cd.rows = NewRows(newName, raw.IDName(), raw.DB())
	cd.deleted = cd.rows.ToDeletedRows()

	return cd
}

// Rows returns the rows being cleaned.
func (cd *CleanRows) Rows() *Rows {
	return cd.rows
}

// DeletedRows returns the rows that were removed during cleaning.
func (cd *CleanRows) DeletedRows() *Rows {
	return cd.deleted
}

// Source returns the source raw data entity.
func (cd *CleanRows) Source() etl.Table {
	return cd.raw
}

// Data returns the cleaned data table (rows).
func (cd *CleanRows) Data() etl.Table {
	return cd.rows
}

// Deleted returns the table containing rows that were removed during cleaning.
func (cd *CleanRows) Deleted() etl.Table {
	return cd.deleted
}
