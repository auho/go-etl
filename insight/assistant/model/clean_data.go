package model

import (
	"github.com/auho/go-etl/v3/job"
)

var _ job.CleanResource = (*CleanData)(nil)

type CleanData struct {
	model

	rows    *Rows
	data    *Data
	deleted *Rows
}

func NewCleanData(rows *Rows) *CleanData {
	cd := &CleanData{}
	cd.rows = rows
	cd.data = rows.ToData()
	cd.deleted = rows.ToDeletedRows()

	return cd
}

func (cd *CleanData) Source() job.Table {
	return cd.rows
}

func (cd *CleanData) Data() job.Table {
	return cd.data
}

func (cd *CleanData) Deleted() job.Table {
	return cd.deleted
}
