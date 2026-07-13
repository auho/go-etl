package entity

import (
	"github.com/auho/go-etl/v3/task/pipeline/database/etl"
)

var _ etl.CleanResource = (*CleanData)(nil)

// CleanData wraps the data cleaning lifecycle for a Data entity.
// It provides access to the source (Rows), the cleaned data (Data), and the deleted rows (Rows).
// This is used in the ETL pipeline's clean stage where data is filtered and validated.
type CleanData struct {
	base // embedded base for command hook and DB connection

	rows    *Rows // source rows before cleaning
	data    *Data // cleaned data output
	deleted *Rows // rows that were removed during cleaning
}

// NewCleanData creates a new CleanData from the given Rows entity.
// It automatically derives the Data and deleted Rows from the source rows.
func NewCleanData(rows *Rows) *CleanData {
	cd := &CleanData{}
	cd.rows = rows
	cd.data = rows.ToData()
	cd.deleted = rows.ToDeletedRows()

	return cd
}

// Source returns the source table (rows before cleaning).
func (cd *CleanData) Source() etl.Table {
	return cd.rows
}

// Data returns the cleaned data table.
func (cd *CleanData) Data() etl.Table {
	return cd.data
}

// Deleted returns the table containing rows that were removed during cleaning.
func (cd *CleanData) Deleted() etl.Table {
	return cd.deleted
}
