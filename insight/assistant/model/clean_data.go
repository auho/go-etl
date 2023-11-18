package model

import (
	"github.com/auho/go-etl/v2/job"
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

func (cd *CleanData) SourceTarget() job.Target {
	return cd.rows
}

func (cd *CleanData) DataTarget() job.Target {
	return cd.data
}

func (cd *CleanData) DeletedTarget() job.Target {
	return cd.deleted
}
